package collect_test

import (
	"io/ioutil"

	"github.com/gocolly/colly"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sanity-io/litter"

	. "github.com/markelog/map/collect"
	"github.com/markelog/map/io"
)

var _ = Describe("io", func() {
	var (
		data *Collect
	)

	BeforeEach(func() {
		html, _ := ioutil.ReadFile("testdata/test.html")
		doc, _ := io.MakeDoc(html)

		data = New(doc)
	})

	AfterEach(func() {
		data = nil
	})

	Describe("Title", func() {
		It("Gets title", func() {
			Expect(data.Title()).To(Equal("test"))
		})
	})

	Describe("Images", func() {
		It("Gets images", func() {
			Expect(data.Images()[0]).To(Equal("test.png"))
		})
	})

	Describe("Styles", func() {
		It("Gets images", func() {
			Expect(data.Styles()[0]).To(Equal("test.css"))
		})
	})

	Describe("Scripts", func() {
		It("Gets scripts", func() {
			Expect(data.Scripts()[0]).To(Equal("test.js"))
		})
	})

	Describe("Video", func() {
		It("Gets video", func() {
			result := data.Video()

			Expect(result).To(ContainElement("test.gif"))
			Expect(result).To(ContainElement("test.webm"))
			Expect(result).To(ContainElement("test.mp4"))
			Expect(result).To(ContainElement("test.ogv"))

			Expect(len(result)).To(Equal(4))
		})
	})

	Describe("Audio", func() {
		It("Gets audio", func() {
			result := data.Video()

			Expect(result).To(ContainElement("test.gif"))
			Expect(result).To(ContainElement("test.webm"))
			Expect(result).To(ContainElement("test.mp4"))
			Expect(result).To(ContainElement("test.ogv"))

			Expect(len(result)).To(Equal(4))
		})
	})

	Describe("Links", func() {
		It("Gets links", func() {
			request := &colly.Request{}
			Expect(data.Links(request)[0]).To(Equal("https://github.com"))
		})
	})

	Describe("Assets", func() {
		It("Gets assets", func() {
			expected := `map[string][]string{
  "audio": []string{
    "foo.wav",
  },
  "images": []string{
    "test.png",
  },
  "scripts": []string{
    "test.js",
  },
  "styles": []string{
    "test.css",
  },
  "video": []string{
    "test.gif",
    "test.webm",
    "test.mp4",
    "test.ogv",
  },
}`

			Expect(litter.Sdump(data.Assets())).To(Equal(expected))
		})
	})
})
