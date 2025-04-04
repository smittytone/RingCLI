package ringcliUtils

import (
	"os"
	"path/filepath"
	errors "ringcli/lib/errors"
	log "ringcli/lib/log"
	"time"
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
		log.ReportError("Could not read bound ring address")
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
		log.ReportErrorAndExit(errors.ERROR_CODE_BINDING_FILE_ERROR, "Could not create the binding store directory (%s)", err)
	}

	// Check if the store itself is present
	bindingStore := filepath.Join(bindingStoreDirectory, "binding")
	_, err = os.Stat(bindingStore)
	if err == nil && !overwrite {
		// The store is present but the user has not marked the operation with `--overwrite`
		log.ReportErrorAndExit(errors.ERROR_CODE_BINDING_FILE_ERROR, "Binding already present. Use --overwrite to replace it")
	}

	// Write out a fresh or replacement binding
	fileData := []byte(address)
	err = os.WriteFile(bindingStore, fileData, 0644)
	if err != nil {
		log.ReportErrorAndExit(errors.ERROR_CODE_BINDING_FILE_ERROR, "Could not store binding")
	}
}

func ToBCD(data int) (byte, error) {

	if data > 99 || data < 0 {
		return 0, &errors.RingcliError{
			Message: "Unsuitable value for BCD conversion",
			Code:    errors.ERROR_CODE_BAD_BCD_INPUT_VALUE,
		}
	}

	return byte(((data / 10) << 4) | (data % 10)), nil
}

func StringifyMonth(month int) string {

	months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	return months[month-1]
}
