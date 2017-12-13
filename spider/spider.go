// Package spider goes through the site tree
// and outputs the meta tree data consumable for reporters
package spider

import (
	"net/url"
	"strings"
	"sync"

	"github.com/gocolly/colly"

	"github.com/markelog/map/collect"
	"github.com/markelog/map/io"
	"github.com/markelog/map/list"
	"github.com/markelog/map/validation"
)

// Result spider data that we eventually return
type Result struct {
	Assets   map[string][]string `json:"assets"`
	URL      string              `json:"url"`
	Name     string              `json:"name"`
	Links    []string            `json:"links"`
	Broken   []string            `json:"broken"`
	Children []*Result           `json:"children"`
	parent   *Result
}

// Progress intermediate data
type Progress struct {
	Data  *Result
	Done  bool
	Error error
}

// Spider configuration
type Spider struct {
	Result *Result
	Error  error

	Progress chan *Progress
	isDone   bool

	waitGroup *sync.WaitGroup
	mutex     *sync.Mutex
	list      *list.List

	path       string
	collector  *colly.Collector
	validation *validation.Validation
}

// New returns new instance of Spider
func New(path, domains string) *Spider {
	var (
		data, _        = url.Parse(path)
		domain         = data.Host
		allowedDomains = []string{domain}
	)

	if len(domains) > 0 {
		allowedDomains = append(strings.Split(domains, ","), domain)
	}

	collector := colly.NewCollector()
	collector.AllowedDomains = allowedDomains

	// Be explicit
	collector.AllowURLRevisit = false

	return &Spider{
		isDone:   false,
		Progress: make(chan *Progress),

		waitGroup: &sync.WaitGroup{},
		mutex:     &sync.Mutex{},
		list:      list.New(),

		path:       path,
		collector:  collector,
		validation: validation.New(path),
	}
}

// Crawl visits the sites and collects data from them
func (spider *Spider) Crawl() (progress chan *Progress) {
	spider.setError()
	spider.setWalker()

	spider.collector.Visit(spider.path)

	go func() {
		spider.waitGroup.Wait()

		spider.mutex.Lock()
		spider.isDone = true
		close(spider.Progress)
		spider.mutex.Unlock()
	}()

	return spider.Progress
}

// Get final result
func (spider Spider) Get() (*Result, error) {
	return spider.Result, spider.Error
}

// Validate spider input
func (spider *Spider) Validate() (err error) {
	err = spider.validation.Parse()
	if err != nil {
		return
	}

	return spider.validation.Check()
}

// emitData emits data for the intermediate data
func (spider *Spider) emitData(data *Progress) {
	spider.mutex.Lock()
	defer spider.mutex.Unlock()

	if spider.isDone == false {
		spider.Progress <- data
	}
}

// setError sets error handler for the spider
func (spider *Spider) setError() {
	spider.collector.OnError(func(response *colly.Response, err error) {
		spider.mutex.Lock()
		defer spider.mutex.Unlock()

		// If first urls breaks
		if spider.Result == nil {
			spider.waitGroup.Add(1)

			go func() {
				spider.emitData(&Progress{
					Error: err,
				})
				spider.waitGroup.Done()
			}()

			return
		}

		if response.Ctx == nil {
			return
		}

		parent := getParent(response)
		if parent == nil {
			return
		}

		parent.Broken = append(parent.Broken, response.Request.URL.String())
	})
}

// setWalker sets crawler walker
func (spider *Spider) setWalker() {
	spider.collector.OnResponse(func(response *colly.Response) {
		body := response.Body

		// Links might lead to the same page, which we might already
		// tackled, so we have to check the response body instead
		if spider.list.Has(body) {
			return
		}
		spider.list.Add(body)

		doc, err := io.MakeDoc(body)
		if err != nil {
			spider.waitGroup.Add(1)
			go func() {
				spider.emitData(&Progress{
					Error: err,
				})
				spider.waitGroup.Done()
			}()
			return
		}

		collection := collect.New(doc)

		output := &Result{
			Assets: collection.Assets(),
			Name:   collection.Title(),
			Links:  collection.Links(response.Request),
			URL:    response.Request.URL.String(),
		}

		if spider.Result == nil {
			spider.Result = output
		}

		spider.waitGroup.Add(1)
		go func() {
			spider.emitData(&Progress{
				Data: output,
			})
			spider.waitGroup.Done()
		}()

		spider.appendToParent(output, response)
		spider.request(output, output.Links)
	})
}

// request multiple links from provided arguments
func (spider Spider) request(output *Result, links []string) {
	spider.waitGroup.Add(len(links))

	for _, link := range links {
		context := colly.NewContext()
		context.Put("parent", output)

		go func(link string, context *colly.Context) {
			spider.collector.Request("GET", link, nil, context, nil)
			spider.waitGroup.Done()
		}(link, context)
	}
}

// getParent gets parent from the context of the response
func getParent(response *colly.Response) (parent *Result) {
	parentInterface := response.Request.Ctx.GetAny("parent")
	if parentInterface != nil {
		parent = parentInterface.(*Result)
		return
	}

	return nil
}

// appendToParent append node to their parent
func (spider Spider) appendToParent(output *Result, response *colly.Response) {
	spider.mutex.Lock()
	defer spider.mutex.Unlock()

	parent := getParent(response)
	if parent == nil {
		return
	}

	children := parent.Children

	for _, element := range children {
		if element.URL == output.URL {
			return
		}
	}

	parent.Children = append(parent.Children, output)

}
