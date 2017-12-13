// Package print provides methods to pretty printing data into the console
package print

import (
	"log"
	"net/url"
	"os"

	"github.com/go-errors/errors"
	"github.com/mgutz/ansi"

	"github.com/markelog/map/print/spinner"
	"github.com/markelog/map/spider"
)

// isDebug checks if we are in debug mode
func isDebug() bool {
	return os.Getenv("MAP_DEBUG") == "true"
}

// ShowError shows the error
func ShowError(err error) {
	if err == nil {
		return
	}

	var (
		stderr = log.New(os.Stderr, "", 0)
		red    = ansi.Color("> ", "red")
	)

	stderr.Println()

	if isDebug() {
		stderr.Println(errors.Wrap(err, 2).ErrorStack())
		return
	}

	_, ok := err.(*url.Error)
	if ok {
		stderr.Printf(red+"Error on %v:", err)
		return
	}
	stderr.Printf(red+"%v", err)
}

// Error shows the error and exit the program with provided exit code
func Error(err error, exitCode int) {
	if err == nil {
		return
	}

	ShowError(err)

	os.Exit(exitCode)
}

// Spin shows the spinner with additional message
// and returns the exitCode if error occuried
func Spin(progress chan *spider.Progress) (exitCode int) {
	spin := spinner.New()
	exitCode = 0

	// Work spinner
	spin.Start()
	for result := range progress {
		if result.Error != nil {
			exitCode = 1

			spin.Stop()
			ShowError(result.Error)
			return
		}

		spin.Set(result.Data.URL)
	}

	spin.Stop()

	return
}
