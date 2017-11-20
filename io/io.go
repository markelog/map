package io

import (
	"io/ioutil"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-errors/errors"
	"golang.org/x/net/html"
)

// WriteFile writes file in a bit more convenient way
func WriteFile(path, content string) (err error) {
	data := []byte(content)

	err = ioutil.WriteFile(path, data, 0700)
	if err != nil {
		return errors.New(err)
	}

	return nil
}

// MakeDoc creates a HTML document from the bytes
func MakeDoc(body []byte) (*goquery.Document, error) {
	text := strings.NewReader(string(body))
	dom, err := html.Parse(text)

	if err != nil {
		return nil, errors.New(err)
	}

	return goquery.NewDocumentFromNode(dom), nil
}
