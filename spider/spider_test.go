package spider_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/sanity-io/litter"

	. "github.com/markelog/map/spider"
)

var _ = Describe("spider", func() {
	var (
		ts      *httptest.Server
		spidy   *Spider
		html, _ = ioutil.ReadFile("testdata/test.html")
	)

	BeforeEach(func() {
		ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			io.WriteString(w, string(html))
		}))

		spidy = New(ts.URL, 20)
	})

	AfterEach(func() {
		ts.Close()
		ts = nil
	})

	Describe("Validate", func() {
		It("Should correct validate the input", func() {
			Expect(New("http://github.com", 20).Validate()).ToNot(HaveOccurred())
		})

		It("Should flag invalid input", func() {
			Expect(New("github.com:20", 20).Validate()).To(HaveOccurred())
		})
	})

	Describe("Crawl", func() {
		It("Should correct validate the input", func() {
			var (
				data, _    = ioutil.ReadFile("testdata/data.txt")
				json       = string(data)
				urlData, _ = url.Parse(ts.URL)
				expected   = strings.Replace(json, "$URL", urlData.Host, 1)
				progress   = spidy.Crawl()
			)

			expected = strings.TrimSpace(expected)

			for value := range progress {
				Expect(value.Error).To(BeNil())
				Expect(value.Data).ToNot(BeNil())

				Expect(litter.Sdump(value.Data)).To(Equal(expected))
			}
		})
	})

	Describe("Get", func() {
		It("Should correct validate the input", func() {
			result, err := spidy.Get()

			Expect(err).ToNot(HaveOccurred())
			Expect(reflect.TypeOf(result).String()).To(Equal("*spider.Result"))
		})
	})
})
