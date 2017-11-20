package io

import (
	"io/ioutil"

	"github.com/go-errors/errors"
)

// WriteFile writes file in a bit more convenient way
func WriteFile(path, content string) error {
	data := []byte(content)

	err := ioutil.WriteFile(path, data, 0700)
	if err != nil {
		return errors.New(err)
	}

	return nil
}
