package rcUtilsCommands

import (
	//"github.com/spf13/cobra"
	// App
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcUtils "ringcli/lib/utils"
)

const (
	ADDRESS_COMMAND_TEXT = "The rings's BLE address. Required"
)

// Globals for all commands
var (
	debug       bool   = false
	ringName    string = ""
	ringAddress string = ""
)

// Set up the `utils` sub-commands' flags.
func init() {
	// Add required flags: --address
	InfoCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)
	//InfoCmd.MarkFlagRequired("address")

	// Add optional flags: --first
	ScanCmd.Flags().BoolVarP(&scanForFirstRing, "first", "f", false, "Stop scanning once first ring found")

	// Add required flags:
	// Add optional flags: --address, --continuous
	FindCmd.Flags().StringVarP(&ringAddress, "address", "", "", ADDRESS_COMMAND_TEXT)
	FindCmd.Flags().BoolVarP(&continuousFlash, "continuous", "c", false, "Flash the ring's LED continuously until cancelled")

	// Add optiona; flags: --address
	ShutdownCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)

	// Add required flags:
	// Add optional flags: --enable, --disable, --period, --address
	SetHeartRateCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)
	SetHeartRateCmd.Flags().BoolVarP(&heartRateEnableSet, "enable", "", false, "Enable periodic heart rate readings")
	SetHeartRateCmd.Flags().BoolVarP(&heartRateDisableSet, "disable", "", false, "Disable periodic heart rate readings")
	SetHeartRateCmd.Flags().IntVarP(&heartRatePeriod, "period", "p", 60, "Heart rate period readings in seconds (0-255)")

	// Add optional flags: --address
	GetHeartRateCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)

	// Add optional flags: --address, --overwrite, --show
	BindCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)
	BindCmd.Flags().BoolVarP(&doOverwrite, "overwrite", "o", false, "Replace an existing binding, if present")
	BindCmd.Flags().BoolVarP(&doShow, "show", "s", false, "Show an existing binding, if present")
}

func getRingAddress() {

	// Check that `ringAddress` has been set by option: if it has, `ringAddress` will not be empty
	if ringAddress == "" {
		// Try to get a stored (bound)
		ringAddress = rcUtils.GetStoredRingAddress()

		if ringAddress == "" {
			// No loaded address so report and bail
			rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
		}
	}
}
