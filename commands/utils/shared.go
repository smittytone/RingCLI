package rcUtilsCommands

import (
	//"github.com/spf13/cobra"
)

// Globals for state
var (
	debug            bool = false
	ringName         string = ""
	ringAddress      string = ""
	scanForFirstRing bool = false
)

// Set up the `utils` sub-commands' flags.
func init() {
	// Add required flags: --address
	InfoCmd.Flags().StringVarP(&ringAddress, "address", "", "", "The rings's BLE address. Required")
	InfoCmd.MarkFlagRequired("address")

	// Add optional flags: --first
	ScanCmd.Flags().BoolVarP(&scanForFirstRing, "first", "f", false, "Stop scanning once first ring found")
}
