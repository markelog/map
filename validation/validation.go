// Package validation provides simple methods for URL validation
package validation

import (
	"net/url"

	"github.com/go-errors/errors"
)

// Validation struct instance
type Validation struct {
	URL    string
	parsed *url.URL
}

// New returns new instance of the Validation
func New(URL string) *Validation {
	return &Validation{
		URL:    URL,
		parsed: nil,
	}
}

// Parse parses URL
func (validation *Validation) Parse() (err error) {
	url, err := url.Parse(validation.URL)
	if err != nil {
		return
	}

	validation.parsed = url

	return
}

// Check execute all the checks
func (validation Validation) Check() (err error) {
	err = validation.CheckScheme()
	if err != nil {
		return
	}

	err = validation.CheckHost()
	if err != nil {
		return
	}

	return
}

// CheckScheme checks the scheme
func (validation Validation) CheckScheme() error {
	if validation.parsed.Scheme == "" {
		return errors.New("You need to specify the protocol")
	}

	return nil
}

// CheckHost checks the host
func (validation Validation) CheckHost() error {
	if validation.parsed.Host == "" {
		return errors.New("You need to specify the host")
	}

	return nil
}
