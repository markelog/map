package json_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/map/reporters/json"
	"github.com/markelog/map/spider"
)

var _ = Describe("reporters", func() {
	Describe("JSON", func() {
		It("Executes json reporter", func() {
			expected := `{"assets":null,"url":"","name":"test","links":null,"broken":null,"children":null}`
			result, _ := Execute(&spider.Result{Name: "test"})

			Expect(result).To(Equal(expected))
		})
	})
})
