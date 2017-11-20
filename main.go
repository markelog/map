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

// Max limits the recursion depth of visited URLs
var max int

// Out is the path data file
var out string

// Command example
const example = `
  Create map and output it to the terminal1
  $ map http://example.com

  Create map and output map in yaml form
  $ map http://example.com --reporter=yaml

  Specify maximum amount of links to check
  $ map http://example.com -r yaml --max=50

	Pipe it
  $ map http://example.com -r yaml -m 50 > example.com.json

	Or use "out" flag (so you can see the spinner :)
  $ map http://example.com -r yaml -m 50 --out=./example.com.json
`

// Command config
var Command = &cobra.Command{
	Use:     "map https://example.com",
	Short:   "Create map of the site",
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
		err := errors.New(`Repoter "` + reporter + `" doesn't not exist`)
		print.Error(err, 2)

		return
	}

	crawler := spider.New(args[0], max)

	// Validate the input
	print.Error(crawler.Validate(), 2)

	// Crawl the site, show the spinner and determine the exit code
	exitCode := print.Spin(crawler.Crawl())

	// Get the result and send it to the reporter
	data, err := crawler.Get()
	print.Error(err, 1)

	serialized, err := reporters.Execute(reporter, data)
	print.Error(err, 1)

	// Either print to the console or save it to a file
	if out == "" {
		fmt.Println(serialized)
	} else {
		print.Error(io.WriteFile(out, serialized), 1)
	}

	os.Exit(exitCode)
}

// Init
func init() {
	cobra.OnInitialize()

	flags := Command.PersistentFlags()

	flags.IntVarP(
		&max,
		"max",
		"m",
		0,
		"Max limits the recursion depth of visited URLs",
	)

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
}

// Main
func main() {
	Command.Execute()
}
