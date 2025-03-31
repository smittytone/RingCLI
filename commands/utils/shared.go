package rcUtilsCommands

import (
	//"github.com/spf13/cobra"
)

const (
	ADDRESS_COMMAND_TEXT = "The rings's BLE address. Required"
)

// Globals for state
var (
	debug            bool = false
	ringName         string = ""
	ringAddress      string = ""
	scanForFirstRing bool = false
	continuousFlash  bool = false
)

// Set up the `utils` sub-commands' flags.
func init() {
	// Add required flags: --address
	InfoCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)
	InfoCmd.MarkFlagRequired("address")

	// Add optional flags: --first
	ScanCmd.Flags().BoolVarP(&scanForFirstRing, "first", "f", false, "Stop scanning once first ring found")

	// Add required flags: --address
	FindCmd.Flags().StringVarP(&ringAddress, "address", "", "", ADDRESS_COMMAND_TEXT)
	FindCmd.MarkFlagRequired("address")
	// Add optional flags: --continuous
	FindCmd.Flags().BoolVarP(&continuousFlash, "continuous", "c", false, "Flash the ring's LED continuously until cancelled")

	// Add required flags: --address
	ShutdownCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)
	ShutdownCmd.MarkFlagRequired("address")
}
