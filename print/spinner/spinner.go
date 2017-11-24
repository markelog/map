// Package spinner creates a spinner
package spinner

import (
	"fmt"
	"os"
	"sync"
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

	waitGroup *sync.WaitGroup
	mutex     *sync.Mutex
	cursed    *curse.Cursor
	spin      *spin.Spinner
}

// New returns the spin instance
func New() *Spin {
	cursed, _ := curse.New()
	fd := int(os.Stdout.Fd())

	return &Spin{
		cursed: cursed,
		spin:   spin.New(),

		mutex: &sync.Mutex{},

		// Do not show the spinner if are in the pipe
		isDone: terminal.IsTerminal(fd) == false,
	}
}

// Start the spinner
func (me *Spin) Start() {
	// If we stopped – then do not try to do it again
	if me.isDone {
		return
	}

	// Just Spinner and additional text with it
	fmt.Println()
	fmt.Println()
	fmt.Println()

	go func() {
		for me.isDone == false {
			me.mutex.Lock()

			me.cursed.MoveUp(2)
			me.cursed.EraseCurrentLine()

			fmt.Print(cyan, me.spin.Next(), reset, " ")

			if len(me.Message) != 0 {
				fmt.Print(gray, truncate(me.Message, 60), reset)
			}
			me.mutex.Unlock()

			fmt.Println()
			fmt.Println()

			time.Sleep(timeout)
		}

		me.Stop()
	}()
}

// Set defines the message on the Spin instance
func (me *Spin) Set(message string) {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	if message == me.Message {
		return
	}

	me.Message = message
}

// Stop the spinner
func (me *Spin) Stop() {
	me.mutex.Lock()
	defer me.mutex.Unlock()

	// If we already finished - don't do anything
	if me.isDone {
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

func truncate(str string, num int) string {
	bnoden := str
	if len(str) > num {
		if num > 3 {
			num -= 3
		}
		bnoden = str[0:num] + "..."
	}
	return bnoden
}
