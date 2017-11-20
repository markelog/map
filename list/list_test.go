package list_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/map/list"
)

var _ = Describe("list", func() {
	It("Correctly check presence of the value", func() {
		test := New()
		value := []byte("test")

		test.Add(value)

		Expect(test.Has(value)).To(Equal(true))
	})

	It("Correctly check absence of the value", func() {
		test := New()
		value := []byte("test")

		test.Add(value)

		Expect(test.Has([]byte("different"))).To(Equal(false))
	})
})
