/*
The MIT License (MIT)

Copyright (c) 2016 Lau, Chok Sheak (for software "timeit" and "timeme")

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"time"
)

/**************************************************************************/

// Constants.

const (
	printPrefix = "[timeme] "
)

/**************************************************************************/

// Global variables.

var (
	isWindows = (os.PathSeparator == '\\') && (os.PathListSeparator == ';')
	newline   = whichString(isWindows, "\r\n", "\n")
)

/**************************************************************************/

// Utilities.

func whichString(result bool, ifTrue, ifFalse string) string {
	if result {
		return ifTrue
	}
	return ifFalse
}

func atoi(s, name string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		fmt.Printf("%v%v \"%v\" not a number: %v%v", printPrefix, name, s, err, newline)
		os.Exit(1)
	}
	return i
}

/**************************************************************************/

// Print and ignore all signals.

func setupSignalHandler() {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel)

	go func() {
		for {
			caughtSignal := <-channel

			fmt.Printf("%vIgnoring signal %v (%#v)%v",
				printPrefix,
				caughtSignal,
				caughtSignal,
				newline)
		}
	}()
}

/**************************************************************************/

// Parse arguments.

func parseArguments(arguments []string) (exitCode int, sleepMilliseconds int) {
	// Print usage.
	if len(arguments) != 3 {
		fmt.Printf("Usage: timeme [exitcode] [sleep-milliseconds]%v", newline)
		os.Exit(1)
	}
	exitCodeString, sleepMillisecondsString := arguments[1], arguments[2]

	// Shortcut mode.
	if exitCodeString == "0" && sleepMillisecondsString == "0" {
		os.Exit(0)
	}

	// Read arguments.
	exitCode = atoi(exitCodeString, "exitcode")
	sleepMilliseconds = atoi(sleepMillisecondsString, "sleep-milliseconds")
	return
}

/**************************************************************************/

// Go!.

func main() {
	// Read arguments.
	exitCode, sleepMilliseconds := parseArguments(os.Args)

	// Catch all incoming signals.
	setupSignalHandler()

	// Generate output for both stdout and stderr.
	fmt.Printf("%vThis string goes to stdout.%v", printPrefix, newline)
	fmt.Fprintf(os.Stderr, "%vThis string goes to stderr.%v", printPrefix, newline)

	// Sleep.
	if sleepMilliseconds > 0 {
		fmt.Printf("%vSleeping for %v milliseconds ...%v", printPrefix, sleepMilliseconds, newline)
		duration, _ := time.ParseDuration(strconv.Itoa(sleepMilliseconds) + "ms")
		time.Sleep(duration)
	}

	// Exit.
	fmt.Printf("%vExiting with code %v%v", printPrefix, exitCode, newline)
	os.Exit(exitCode)
}

/**************************************************************************/
