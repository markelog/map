package validation_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/markelog/map/validation"
)

var _ = Describe("validation", func() {
	Describe("Parse", func() {
		It("Correctly parses url", func() {
			Expect(
				New("http://github.com").Parse(),
			).NotTo(HaveOccurred())
		})

		It("Returns error when url is incorrect", func() {
			Expect(
				New("@#$%^&*()").Parse(),
			).To(HaveOccurred())
		})
	})

	Describe("CheckHost", func() {
		It("Correctly checks host", func() {
			validation := New("http://github.com")
			validation.Parse()

			Expect(
				validation.CheckHost(),
			).NotTo(HaveOccurred())
		})

		It("Correctly checks host", func() {
			validation := New("http://")
			validation.Parse()

			Expect(
				validation.CheckHost(),
			).To(HaveOccurred())
		})
	})

	Describe("CheckScheme", func() {
		It("Correctly checks host", func() {
			validation := New("http://github.com")
			validation.Parse()
			Expect(
				validation.CheckScheme(),
			).NotTo(HaveOccurred())
		})

		It("Returns error when url is incorrect", func() {
			validation := New("github.com")
			validation.Parse()
			Expect(
				validation.CheckScheme(),
			).To(HaveOccurred())
		})
	})

	Describe("Check", func() {
		It("Correctly checks all", func() {
			validation := New("http://github.com")
			validation.Parse()

			Expect(
				validation.Check(),
			).NotTo(HaveOccurred())
		})

		It("Returns checks should not pass", func() {
			validation := New("github.com")
			validation.Parse()

			Expect(
				validation.Check(),
			).To(HaveOccurred())
		})
	})
})
