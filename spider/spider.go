package spider

import (
	"net/url"
	"runtime"

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
	Error error
}

// Spider configuration
type Spider struct {
	Result   *Result
	Progress chan *Progress
	Error    error

	isDone     bool
	domain     string
	list       *list.List
	collector  *colly.Collector
	validation *validation.Validation
}

// New returns new instance of Spider
func New(domain string, max int) *Spider {
	data, _ := url.Parse(domain)

	collector := colly.NewCollector()
	collector.AllowedDomains = []string{data.Host}
	collector.MaxDepth = max

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: runtime.NumCPU(),
	})

	return &Spider{
		Progress: make(chan *Progress),

		list:       list.New(),
		collector:  collector,
		validation: validation.New(domain),
		domain:     domain,
		isDone:     false,
	}
}

// Crawl visits the sites and collects data from them
func (spider *Spider) Crawl() (progress chan *Progress) {
	spider.setError()
	spider.setWalker()

	go func() {
		spider.Error = spider.collector.Visit(spider.domain)

		spider.collector.Wait()

		spider.isDone = true
		close(spider.Progress)
	}()

	return spider.Progress
}

// Get final result
func (spider Spider) Get() (*Result, error) {
	spider.collector.Wait()

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
	if spider.isDone == false {
		spider.Progress <- data
	}
}

// setError sets error handler for the spider
func (spider *Spider) setError() {
	spider.collector.OnError(func(response *colly.Response, err error) {
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
			go spider.emitData(&Progress{
				Error: err,
			})
			return
		}

		collection := collect.New(doc)

		output := &Result{
			Assets: collection.Assets(),
			Name:   collection.Title(),
			Links:  collection.Links(response.Request),
			URL:    response.Request.URL.String(),
		}

		if getParent(response) == nil {
			spider.Result = output
		}

		go spider.emitData(&Progress{
			Data: output,
		})

		spider.appendToParent(output, response)
		spider.request(output, output.Links)
	})
}

// request multiple links from provided arguments
func (spider Spider) request(output *Result, links []string) {
	for _, link := range links {
		context := colly.NewContext()
		context.Put("parent", output)

		go spider.collector.Request("GET", link, nil, context, nil)
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
	parent := getParent(response)
	if parent != nil {
		parent.Children = append(parent.Children, output)
	}
}
