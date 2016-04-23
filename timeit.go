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
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

/**************************************************************************/

// Constants.

const (
	version = "0.1.20160422"

	// Thousand.
	exp3 = 1000

	// Million.
	exp6 = 1000000

	// Billion.
	exp9 = 1000000000

	// Maximum value to display in nanoseconds.
	nanoDisplayMax = 99999

	// Standardize timeit output with a prefix.
	printPrefix = "[timeit] "
)

/**************************************************************************/

// Global variables.

var (
	isWindows = (os.PathSeparator == '\\') && (os.PathListSeparator == ';')
	newline   = whichString(isWindows, "\r\n", "\n")

	// Allows injection of custom time (in nanoseconds) to print as output.
	testDuration = time.Duration(-1)

	startTime        time.Time
	endTime          time.Time
	commandToExecute *exec.Cmd
)

/**************************************************************************/

// Utilities.

func whichString(result bool, ifTrue, ifFalse string) string {
	if result {
		return ifTrue
	}
	return ifFalse
}

func ps(s string) {
	fmt.Print(s)
}

func pl() {
	ps(newline)
}

func pf(format string, a ...interface{}) {
	// Make sure newline does not break up from the rest of the string.
	// This could happen when user presses Ctrl+C.
	ps(fmt.Sprintf(format, a...) + newline)
}

func getExitCode(cmd *exec.Cmd, err error) (int, bool) {
	var waitStatus syscall.WaitStatus
	if err != nil {
		exitError, isRightType := err.(*exec.ExitError)
		if isRightType {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
		} else {
			return 0, false
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
	}
	return waitStatus.ExitStatus(), true
}

// Basically append leading zeros.
func getAsDecimals(number int64, numDigits int) string {
	decimalsAsRunes := make([]rune, numDigits)
	numberAsRunes := []rune(strconv.FormatInt(number, 10))
	j := len(numberAsRunes)
	for i := numDigits - 1; i >= 0; i-- {
		j--
		if j >= 0 {
			decimalsAsRunes[i] = numberAsRunes[j]
		} else {
			decimalsAsRunes[i] = '0'
		}
	}
	return string(decimalsAsRunes)
}

/**************************************************************************/

// Making the duration readable.

func writeReadableTimeUnit(buffer *bytes.Buffer, numUnits int64, singularString, pluralString string) {
	if numUnits <= 0 {
		return
	}
	if buffer.Len() != 0 {
		buffer.WriteRune(' ')
	}
	buffer.WriteString(strconv.FormatInt(numUnits, 10))
	if numUnits > 1 {
		buffer.WriteString(pluralString)
	} else {
		buffer.WriteString(singularString)
	}
}

// This string is meant to be as human-readable as possible.
func getDurationAsReadableString(duration time.Duration) string {
	if duration == time.Duration(0) {
		return "0ns"
	}

	nanos := duration.Nanoseconds()

	seconds := nanos / exp9
	nanos %= exp9

	minutes := seconds / 60
	seconds %= 60

	hours := minutes / 60
	minutes %= 60

	days := hours / 24
	hours %= 24

	var buffer bytes.Buffer

	writeReadableTimeUnit(&buffer, days, " day", " days")
	writeReadableTimeUnit(&buffer, hours, " hour", " hours")
	writeReadableTimeUnit(&buffer, minutes, " minute", " minutes")

	// Write seconds and under.
	nanosRemaining := (seconds * exp9) + nanos
	remainingDuration := time.Duration(nanosRemaining)

	if buffer.Len() != 0 {
		buffer.WriteRune(' ')
	}
	buffer.WriteString(getDurationAsSimpleString(remainingDuration))

	return buffer.String()
}

// This string is meant to be easy to copy and paste for use in computations.
func getDurationAsSimpleString(duration time.Duration) string {
	nanos := duration.Nanoseconds()

	if nanos <= nanoDisplayMax {
		return strconv.FormatInt(nanos, 10) + "ns"
	}

	if nanos >= exp9 {
		// Return in seconds.
		seconds := nanos / exp9
		nanos %= exp9
		millis := nanos / exp6
		return strconv.FormatInt(seconds, 10) + "." + getAsDecimals(millis, 3) + "s"
	}

	// Return in milliseconds.
	millis := nanos / exp6
	nanos %= exp6
	micros := nanos / exp3
	return strconv.FormatInt(millis, 10) + "." + getAsDecimals(micros, 3) + "ms"
}

/**************************************************************************/

// Main flow logic of this program.

func execAndTime(command string, args ...string) {
	// Create command.
	commandToExecute = exec.Command(command, args...)
	commandToExecute.Stdout = os.Stdout
	commandToExecute.Stderr = os.Stderr

	// Execute with timing.
	startTime = time.Now()
	err := commandToExecute.Run()
	endTime = time.Now()

	// Get exit code.
	exitCode, hasExitCode := getExitCode(commandToExecute, err)

	// Blank line to mark start of timeit output.
	pl()

	// Print error if any.
	// No need to print out non-zero exit status (noise).
	if err != nil && !strings.HasPrefix(err.Error(), "exit status ") {
		pf("%vError: %v", printPrefix, err)
	}

	// Exit now if command was not executed at all.
	if !hasExitCode {
		os.Exit(-1)
	}

	// Show how long it took.
	printTiming()

	// Exit with the same exit code as the given command.
	os.Exit(exitCode)
}

func printTiming() {
	// Get durations.
	duration := endTime.Sub(startTime)
	systemTime := time.Duration(0)
	userTime := time.Duration(0)

	if commandToExecute != nil && commandToExecute.ProcessState != nil {
		systemTime = commandToExecute.ProcessState.SystemTime()
		userTime = commandToExecute.ProcessState.UserTime()
	}

	// Inject given duration for testing.
	if testDuration >= 0 {
		duration = testDuration
		systemTime = testDuration
		userTime = testDuration
	}

	// Print final output.
	var buffer bytes.Buffer

	buffer.WriteString(printPrefix)

	durationLongString := getDurationAsReadableString(duration)
	durationShortString := getDurationAsSimpleString(duration)

	if durationLongString == durationShortString {
		buffer.WriteString("Real ")
		buffer.WriteString(durationShortString)
	} else {
		buffer.WriteString(durationLongString)
		buffer.WriteString(" (")
		buffer.WriteString(durationShortString)
		buffer.WriteString(")")
	}

	buffer.WriteString(", Usr ")
	buffer.WriteString(getDurationAsSimpleString(userTime))

	buffer.WriteString(", Sys ")
	buffer.WriteString(getDurationAsSimpleString(systemTime))

	buffer.WriteString(" (")
	buffer.WriteString(startTime.Format("15:04:05"))
	buffer.WriteString("-")
	buffer.WriteString(endTime.Format("15:04:05"))
	buffer.WriteString(")")

	ps(buffer.String())
	pl()
}

/**************************************************************************/

// Passthrough all signals.

func setupSignalPassThrough() {
	channel := make(chan os.Signal, 1)
	signal.Notify(channel)

	go func() {
		for {
			caughtSignal := <-channel

			if commandToExecute != nil &&
				commandToExecute.Process != nil {
				pf("%vSending %v (%#v) to command", printPrefix, caughtSignal, caughtSignal)
				commandToExecute.Process.Signal(caughtSignal)
			} else {
				endTime = time.Now()
				pl()
				pf("%vGoodbye", printPrefix)
				printTiming()
				os.Exit(-1)
			}
		}
	}()
}

/**************************************************************************/

// Arguments parsing.

// Returns true if s begins with '-' and is followed by digits only.
func isOptionString(s string) bool {
	if len(s) <= 1 {
		return false
	}
	for pos, char := range s {
		if pos == 0 {
			if char != '-' {
				return false
			}
		} else {
			if (char < '0') || (char > '9') {
				return false
			}
		}
	}
	return true
}

func parseArguments(args []string) []string {
	// Allow specification of timeit duration for testing.
	// This feature does not need to be publicly-accessible.
	// No need to hide it because there are no security concerns.
	if (len(args) > 1) && isOptionString(args[1]) {
		s := args[1][1:]
		i, _ := strconv.Atoi(s)
		testDuration = time.Duration(i)
		// Remove first argument from list.
		args = append(args[0:1], args[2:]...)
		pf("%vtestDuration = %v nanoseconds", printPrefix, testDuration)
		pl()
	}

	// Check arguments.
	if len(args) == 1 {
		printUsage()
		os.Exit(1)
	}

	return args[1:]
}

/**************************************************************************/

// Print usage.

func printUsage() {
	pl()
	pf("USAGE: timeit <command> [arg]...")
	pl()
	pf("Copyright (c) under The MIT License for timeit version %v.", version)
	pf("Online: https://github.com/choksheak/timeit/blob/master/LICENSE.txt")
	pl()
}

/**************************************************************************/

// Go!

func main() {
	args := parseArguments(os.Args)
	setupSignalPassThrough()
	execAndTime(args[0], args[1:]...)
}

/**************************************************************************/
