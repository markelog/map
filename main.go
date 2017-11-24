// Package main is the CLI wrapper for the crawler
package main

import (
	"fmt"
	"os"

	"github.com/go-errors/errors"
	"github.com/spf13/cobra"

	"github.com/markelog/map/io"
	"github.com/markelog/map/print"
	"github.com/markelog/map/reporters"
	"github.com/markelog/map/spider"
)

// Reporter name
var reporter string

// Out is the path data file
var out string

// Domains to follow
var domains string

// Command example
const example = `
  Create map and output it to the terminal
  $ map http://example.com

  Create map and output map in yaml form
  $ map http://example.com --reporter=yaml

  Pipe it
  $ map http://example.com -r yaml > example.com.yaml

  Or use "out" flag to pipe (so you can see the spinner comparing with previous command :)
  $ map http://example.com -r yaml --out=./example.com.yaml

	With additional domains
	$ map https://example.com --domains=www.google.ru,www.google.com
`

// Command config
var Command = &cobra.Command{
	Use:     "map https://example.com",
	Short:   "Site mapper",
	Example: example,
	Run:     Run,
}

// Run the command!
func Run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		print.Error(errors.New("Target is not specified"), 2)

		return
	}

	// Reporter has to exist
	if reporters.Exist(reporter) == false {
		err := errors.New(`Reporter "` + reporter + `" does not exist`)
		print.Error(err, 2)

		return
	}

	crawler := spider.New(args[0], domains)

	// Validate the input
	print.Error(crawler.Validate(), 2)

	// Crawl the site, show the spinner and determine the exit code
	exitCode := print.Spin(crawler.Crawl())

	// Get the result and send it to the reporter
	data, err := crawler.Get()
	print.Error(err, 1)

	serialized, err := reporters.Execute(reporter, data)
	print.Error(err, 1)

	if len(serialized) > 0 {
		// Either print to the console or save it to a file
		if len(out) > 0 {
			fmt.Println(serialized)
		} else {
			print.Error(io.WriteFile(out, serialized), 1)
		}
	}

	os.Exit(exitCode)
}

// Init
func init() {
	cobra.OnInitialize()

	flags := Command.PersistentFlags()

	flags.StringVarP(
		&reporter,
		"reporter",
		"r",
		"json",
		"Show data in certain representation",
	)

	flags.StringVarP(
		&out,
		"out",
		"o",
		"",
		"Output data to the file without pipe but with the spinner :)",
	)

	flags.StringVarP(
		&domains,
		"domains",
		"d",
		"",
		"Domains to follow (as addition to the base url), comma as a delimter",
	)
}

// Main
func main() {
	Command.Execute()
}
