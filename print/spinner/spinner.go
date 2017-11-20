package spinner

import (
	"fmt"
	"os"
	"time"

	"github.com/markelog/curse"
	"github.com/mgutz/ansi"
	"github.com/tj/go-spin"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	gray  = ansi.ColorCode("240")
	reset = ansi.ColorCode("reset")
	cyan  = ansi.ColorCode("cyan")

	timeout = 200 * time.Millisecond
)

// Spin settings
type Spin struct {
	Message string

	isDone     bool
	isTerminal bool
	cursed     *curse.Cursor
	spin       *spin.Spinner
}

// New returns the spin instance
func New() *Spin {
	cursed, _ := curse.New()
	fd := int(os.Stdout.Fd())

	return &Spin{
		cursed: cursed,
		spin:   spin.New(),

		// Do not show the spinner if are in the pipe
		isDone: terminal.IsTerminal(fd) == false,
	}
}

// Start starts the spinner
func (me *Spin) Start() {

	// If we stopped – then do not try to do it again
	if me.isDone {
		return
	}

	// Just Spinner and additional text with it
	started := false
	go func() {
		for me.isDone == false {
			if started {
				me.cursed.MoveUp(2)
			}
			started = true

			me.cursed.EraseCurrentLine()
			fmt.Print(cyan, me.spin.Next(), reset, " ")

			if len(me.Message) != 0 {
				fmt.Print(gray, me.Message, reset)
			}
			fmt.Println()
			fmt.Println()

			time.Sleep(timeout)
		}
	}()
}

// Set define message on the Spin instance
func (me *Spin) Set(message string) {
	me.Message = message
}

// Stop stops the instance
func (me Spin) Stop() {

	// If we already finished - don't do anything
	if me.isDone == true {
		return
	}

	me.isDone = true

	// Clean-up
	me.cursed.MoveUp(1)
	me.cursed.EraseCurrentLine()
	me.cursed.MoveUp(1)
	me.cursed.EraseCurrentLine()
	me.cursed.MoveUp(1)
	me.cursed.EraseCurrentLine()
}
