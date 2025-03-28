package rcUtilsCommands

import (


	//"github.com/spf13/cobra"
)

// Globals to hold the device ID etc.
var (
	debug            bool = false
	ringName         string = ""
	ringAddress      string = ""
	scanForFirstRing bool = false
)

const (
	DEVICE_INFO_SERVICE_ID                       uint16 = 0x180A
	DEVICE_INFO_SERVICE_MANUFACTURER_CHAR_ID     uint16 = 0x2A29
	DEVICE_INFO_SERVICE_FIRMWARE_VERSION_CHAR_ID uint16 = 0x2A26
	DEVICE_INFO_SERVICE_HARDWARE_VERSION_CHAR_ID uint16 = 0x2A27
	DEVICE_INFO_SERVICE_SYSTEM_ID_CHAR_ID        uint16 = 0x2A23
	DEVICE_INFO_SERVICE_PNP_ID_CHAR_ID           uint16 = 0x2A50
)

// Set up the `utils` sub-commands' flags.
func init() {
	// Add required flags: --address
	InfoCmd.Flags().StringVarP(&ringAddress, "address", "", "", "The rings's BLE address. Required")
	InfoCmd.MarkFlagRequired("address")

	// Add optional flags: --first
	ScanCmd.Flags().BoolVarP(&scanForFirstRing, "first", "f", false, "Stop scanning once first ring found")
}
