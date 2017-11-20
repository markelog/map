package print

import (
	"fmt"
	"os"

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

	fmt.Println()
	fmt.Print(ansi.Color("> ", "red"))

	fmt.Fprintf(os.Stderr, "%v", err)
}

// Error shows the error and exit the program with provided exit code
func Error(err error, exitCode int) {
	if err == nil {
		return
	}

	ShowError(err)

	fmt.Println()
	fmt.Println()
	fmt.Println(`Maybe try "--help" flag?`)
	fmt.Println()

	os.Exit(exitCode)
}

// Spin shows the spinner with additional message
// and returns the exitCode if error occuried
func Spin(progress chan *spider.Progress) (exitCode int) {
	spin := spinner.New()
	exitCode = 0

	// Work spinner
	spin.Start()
	fmt.Println()
	for result := range progress {
		if result.Error != nil {
			exitCode = 2

			ShowError(result.Error)
			fmt.Println()

			continue
		}

		spin.Set(result.Data.URL)
	}
	spin.Stop()

	return
}
