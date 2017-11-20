package io_test

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/bouk/monkey"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/map/io"
)

var _ = Describe("io", func() {
	Describe("WriteFile", func() {
		var (
			filename string
			data     []byte
			perm     os.FileMode
			result   error
		)
		BeforeEach(func() {
			monkey.Patch(ioutil.WriteFile, func(
				_filename string,
				_data []byte,
				_perm os.FileMode,
			) error {
				filename = _filename
				data = _data
				perm = _perm

				return result
			})
		})

		AfterEach(func() {
			filename = ""
			data = nil
			perm = os.FileMode(0)

			monkey.Unpatch(ioutil.WriteFile)
		})

		It("Correctly passes the data to WriteFile method", func() {
			var (
				path = "test-path"
				data = "test-data"

				err = WriteFile(path, data)
			)

			Expect(err).To(BeNil())
			Expect(filename).To(Equal(path))
			Expect(string(data)).To(Equal(data))
		})

		It("Correctly returns an error", func() {
			result = errors.New("test")
			err := WriteFile("test-path", "test-data")

			Expect(err.Error()).To(Equal("test"))
		})
	})

	Describe("MakeDoc", func() {
		It("Correctly parses the document", func() {
			var (
				html     = `<!doctype html><meta charset=utf-8><title>short</title>`
				doc, err = MakeDoc([]byte(html))
				title    = doc.Find("title").Text()
			)

			Expect(err).To(BeNil())
			Expect(reflect.TypeOf(doc).String()).To(Equal("*goquery.Document"))
			Expect(title).To(Equal("short"))
		})
	})
})
