// Package collect provides methods for collecting the data
// from the DOM document
package collect

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

// Collect settings
type Collect struct {
	doc *goquery.Document
}

// New returns new Collect instance
func New(doc *goquery.Document) *Collect {
	return &Collect{
		doc: doc,
	}
}

// Title gives text of the first <title>
func (collect Collect) Title() (text string) {
	title := collect.doc.Find("title")

	if len(title.Nodes) == 0 {
		return ""
	}

	return title.Eq(0).Text()
}

// Images returns all <img> "src" attribute
// TODO: Parse srcset attribute as well
func (collect Collect) Images() (images []string) {
	collect.doc.Find("img").Each(func(i int, node *goquery.Selection) {
		src, exist := node.Attr("src")
		if exist == false {
			return
		}

		images = append(images, src)
	})

	return
}

// Styles returns all <link> "href" attribute
func (collect Collect) Styles() (links []string) {
	collect.doc.Find("link").Each(func(i int, node *goquery.Selection) {
		href, exist := node.Attr("href")
		if exist == false {
			return
		}

		links = append(links, href)
	})

	return
}

// Scripts returns all <script> "src" attribute
func (collect Collect) Scripts() (scripts []string) {
	collect.doc.Find("script").Each(func(i int, node *goquery.Selection) {
		src, exist := node.Attr("src")
		if exist == false {
			return
		}

		scripts = append(scripts, src)
	})

	return
}

// Links returns all <a> "href" attribute
func (collect Collect) Links(request *colly.Request) (links []string) {
	// a[href] selector should be slower
	collect.doc.Find("a").Each(func(i int, node *goquery.Selection) {
		href, exist := node.Attr("href")
		if exist == false {
			return
		}

		links = append(links, request.AbsoluteURL(href))
	})

	return
}

// Video returns all <video>, <kind> and <source> "src" attribute
func (collect Collect) Video() (video []string) {
	collect.doc.Find("video").Each(func(i int, node *goquery.Selection) {
		src, exist := node.Attr("src")
		if exist {
			video = append(video, src)
		}

		src, exist = node.Attr("poster")
		if exist {
			video = append(video, src)
		}

		node.Find("kind").Each(func(i int, node *goquery.Selection) {
			src, exist := node.Attr("src")
			if exist == false {
				return
			}

			video = append(video, src)
		})

		node.Find("source").Each(func(i int, node *goquery.Selection) {
			src, exist := node.Attr("src")
			if exist == false {
				return
			}

			video = append(video, src)
		})
	})

	return
}

// Audio returns all <audio> and <source> "src" attribute
func (collect Collect) Audio() (audio []string) {
	collect.doc.Find("audio").Each(func(i int, node *goquery.Selection) {
		src, exist := node.Attr("src")
		if exist {
			audio = append(audio, src)
		}

		node.Find("source").Each(func(i int, node *goquery.Selection) {
			src, exist := node.Attr("src")
			if exist == false {
				return
			}

			audio = append(audio, src)
		})
	})

	return
}

// Assets returns all available and supported by us assets
func (collect Collect) Assets() map[string][]string {
	return map[string][]string{
		"images":  collect.Images(),
		"styles":  collect.Styles(),
		"scripts": collect.Scripts(),
		"video":   collect.Video(),
		"audio":   collect.Audio(),
	}
}
