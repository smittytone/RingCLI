package ringcliLog

import (
	"fmt"
	"os"
	spinner "ringcli/lib/spinner"
	"strings"
)

const (
	PLAIN_MESSAGE   int = 0
	ERROR_MESSAGE   int = 1
	WARNING_MESSAGE int = 2
	DEBUG_MESSAGE   int = 3
	DATA_OUTPUT     int = 4

	ESC    string = "\x1B"
	CSI    string = ESC + "["
	CURSOR string = "|/-\\"
)

var (
	bspCount      int = 0
	cursorSpinner *spinner.Spinner
)

func CursorLeft() {

	raw(CSI + "0G")
}

func CursorHide() {

	raw(CSI + "?25l")
}

func CursorShow() {

	raw(CSI + "?25h")
}

func Prompt(text string) {

	bspCount = raw(text + "...  ")
	cursorSpinner = spinner.NewSpinner(CURSOR)
	cursorSpinner.Start()
}

func RealtimeDataOut(text string) {

	CursorLeft()
	bspCount = raw(text + "  ")
}

func RealtimeDataClear() {

	CursorLeft()
	CursorShow()
	bspCount = 0
}

func ClearPrompt() {

	cursorSpinner.Stop()
	backspaces(bspCount)
	bspCount = 0
}

func raw(msg string, values ...any) int {

	output := fmt.Sprintf(msg, values...)
	fmt.Fprintf(os.Stderr, output)
	return len(output)
}

func backspaces(count int) {

	if count > 0 {
		raw(CSI + fmt.Sprintf("%dD", count))
		raw(strings.Repeat(" ", count))
		raw(CSI + fmt.Sprintf("%dD", count))
	}
}

func backspace(count int) {

	if count > 0 {
		raw(CSI + fmt.Sprintf("%dD", count))
	}
}

func Report(msg string, values ...any) {

	log(PLAIN_MESSAGE, msg, values...)
}

func ReportWarning(errMsg string, values ...any) {

	log(WARNING_MESSAGE, errMsg, values...)
}

func ReportError(errMsg string, values ...any) {

	log(ERROR_MESSAGE, errMsg, values...)
}

/*
func ReportDebug(errMsg string, values ...any) {

	// Only report these messages if the Level is DEBUG or higher
	if mvAppConfig.Config.LogLevel == mvSharedData.LOG_LEVEL_DBG {
		log(DEBUG_MESSAGE, errMsg, values...)
	}
}
*/

func ReportErrorAndExit(errCode int, errMsg string, values ...any) {

	log(ERROR_MESSAGE, errMsg, values...)
	CursorShow()
	os.Exit(errCode)
}

func log(msgType int, msg string, values ...any) {

	if cursorSpinner != nil && cursorSpinner.IsAnimating() {
		cursorSpinner.Stop()
		CursorLeft()
		fmt.Fprintf(os.Stderr, "\n")
	}

	CursorShow()
	outputMsg := msg
	if len(values) > 0 {
		outputMsg = fmt.Sprintf(msg, values...)
	}

	switch msgType {
	case ERROR_MESSAGE:
		fmt.Fprintln(os.Stderr, "[ERROR]", outputMsg)
	case WARNING_MESSAGE:
		fmt.Fprintln(os.Stderr, "[WARNING]", outputMsg)
	case DEBUG_MESSAGE:
		fmt.Fprintln(os.Stderr, "[DEBUG]", outputMsg)
	case DATA_OUTPUT:
		// Date to be output to stdout
		fmt.Fprintln(os.Stdout, outputMsg)
	default:
		// Standard output
		fmt.Fprintln(os.Stderr, outputMsg)
	}
}
