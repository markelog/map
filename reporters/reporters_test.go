package reporters_test

import (
	"io/ioutil"

	"bou.ke/monkey"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/map/reporters"
	"github.com/markelog/map/reporters/json"
	"github.com/markelog/map/spider"
)

var _ = Describe("reporters", func() {
	Describe("Execute", func() {
		var (
			jsonExecuted = false
		)
		BeforeEach(func() {
			monkey.Patch(json.Execute, func(data *spider.Result) (string, error) {
				jsonExecuted = true
				return "", nil
			})
		})

		AfterEach(func() {
			jsonExecuted = false

			monkey.Unpatch(ioutil.WriteFile)
		})

		It("Executes json reporter", func() {
			Execute("json", &spider.Result{})
			Expect(jsonExecuted).To(Equal(true))
		})

		It("Does not executes json reporter", func() {
			Execute("yaml", &spider.Result{})

			Expect(jsonExecuted).To(Equal(false))
		})

		It("Returns error if reporter doesn't exist", func() {
			str, err := Execute("nope", &spider.Result{})

			Expect(str).To(Equal(""))
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("Exist", func() {
		It("Checks if reporter exists", func() {
			Expect(Exist("json")).To(Equal(true))
		})

		It("Checks if reporter does not exists", func() {
			Expect(Exist("nope")).To(Equal(false))
		})
	})
})
