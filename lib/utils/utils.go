package rcUtils

import (
	"os"
	"path/filepath"
	"time"
	// App
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
)

const (
	cursor string = "|/-\\"
)

var (
	animationTimer *time.Ticker
	isAnimating    bool = false
	animationDone  chan bool
	cursorIndex    int = 0
)

func StartToday(date time.Time) time.Time {

	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
}

func GetStoredRingAddress() string {

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		// No home directory!
		return ""
	}

	store := filepath.Join(homeDirectory, ".config", "ringcli", "binding")
	_, err = os.Stat(store)
	if err != nil {
		// No store
		return ""
	}

	bytes, err := os.ReadFile(store)
	if err != nil {
		// Store exists but read failed
		rcLog.ReportError("Could not read bound ring address")
		return ""
	}

	return string(bytes)
}

func createStoreDirectory() (string, error) {

	// Get the user's home directory
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		// No home directory!
		return "", err
	}

	// See if we have a store directory already
	bindingStoreDirectory := filepath.Join(homeDirectory, ".config", "ringcli")
	_, err = os.Stat(bindingStoreDirectory)
	if err == nil {
		// Store already created
		return bindingStoreDirectory, nil
	}

	// Make the store directory
	err = os.MkdirAll(bindingStoreDirectory, os.ModePerm)
	return bindingStoreDirectory, err
}

func MakeBinding(address string, overwrite bool) {

	// Check and make if necessary the store directory
	bindingStoreDirectory, err := createStoreDirectory()
	if err != nil {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BINDING_FILE_ERROR, "Could not create the binding store directory (%s)", err)
	}

	// Check if the store itself is present
	bindingStore := filepath.Join(bindingStoreDirectory, "binding")
	_, err = os.Stat(bindingStore)
	if err == nil && !overwrite {
		// The store is present but the user has not marked the operation with `--overwrite`
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BINDING_FILE_ERROR, "Binding already present. Use --overwrite to replace it")
	}

	// Write out a fresh or replacement binding
	fileData := []byte(address)
	err = os.WriteFile(bindingStore, fileData, 0644)
	if err != nil {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BINDING_FILE_ERROR, "Could not store binding")
	}
}

func ToBCD(data int) byte {

	if data > 99 || data < 0 {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_BCD_INPUT_VALUE, "Unsuitable value for BCD conversion")
	}

	return byte(((data / 10) << 4) | (data % 10))
}

func AnimateCursor() {

	if isAnimating {
		return
	}

	animationTimer = time.NewTicker(50 * time.Millisecond)
	isAnimating = true
	animationDone = make(chan bool)
	go func() {
		for {
			select {
			case <-animationDone:
				//rcLog.Backspaces(1)
				isAnimating = false
				return
			case <-animationTimer.C:
				cursorIndex += 1
				if cursorIndex >= len(cursor) {
					cursorIndex = 0
				}

				rcLog.Backspace(1)
				_ = rcLog.Raw(string(cursor[cursorIndex]))
			}
		}
	}()
}

func StopAnimation() {

	if isAnimating {
		animationDone <- true
	}
}
