package rcUtilsCommands

import (
	errors "ringcli/lib/errors"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
)

const (
	ADDRESS_COMMAND_TEXT = "The rings's BLE address. Required if no ring has been bound."
)

// Globals for all commands
var (
	debug       bool   = false // Are we in debug mode?
	ringAddress string = ""    // Ring BLE address
)

// Set up the `utils` sub-commands' flags.
func init() {
	// BATTERY
	// Add optional flags: --address
	BatteryCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)

	// BIND
	// Add optional flags: --address, --overwrite, --show, --name
	BindCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)
	BindCmd.Flags().StringVarP(&ringName, "name", "n", "not set", "The ring's display name")
	BindCmd.Flags().BoolVarP(&doOverwrite, "overwrite", "o", false, "Replace an existing binding, if present")
	BindCmd.Flags().BoolVarP(&doShow, "show", "s", false, "Show an existing binding, if present")

	// FIND
	// Add optional flags: --address, --continuous
	FindCmd.Flags().StringVarP(&ringAddress, "address", "", "", ADDRESS_COMMAND_TEXT)
	FindCmd.Flags().BoolVarP(&continuousFlash, "continuous", "c", false, "Flash the ring's LED continuously until cancelled")

	// GET HEART RATE
	// Add optional flags: --address, --show, --period, --enable, --disable
	HeartRateCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)
	HeartRateCmd.Flags().BoolVarP(&getHeartRateSettings, "show", "s", false, "Display current periodic heart rate settings")
	HeartRateCmd.Flags().IntVarP(&heartRatePeriod, "period", "p", 60, "Heart rate period readings in seconds (0-255)")
	HeartRateCmd.Flags().BoolVarP(&heartRateEnableSet, "enable", "", false, "Enable periodic heart rate readings")
	HeartRateCmd.Flags().BoolVarP(&heartRateDisableSet, "disable", "", false, "Disable periodic heart rate readings")

	// INFO
	// Add optional flags: --address
	InfoCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)

	// SCAN
	// Add optional flags: --first
	ScanCmd.Flags().BoolVarP(&scanForFirstRing, "first", "f", false, "Stop scanning once first ring found")
	ScanCmd.Flags().BoolVarP(&debug, "all", "", false, "Show all devices")
	ScanCmd.Flags().BoolVarP(&doBind, "bind", "b", false, "Automatically bind a discovered ring")

	// SHUTDOWN
	// Add optional flags: --address
	ShutdownCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)

	// SET TIME
	// Add optional flags: --address
	SetTimeCmd.Flags().StringVarP(&ringAddress, "address", "a", "", ADDRESS_COMMAND_TEXT)
}

func getRingAddress() {

	// Check that `ringAddress` has been set by option: if it has, `ringAddress` will not be empty
	if ringAddress == "" {
		// Try to get a stored (bound)
		ringAddress = utils.GetStoredRingAddress()

		if ringAddress == "" {
			// No loaded address so report and bail
			log.ReportErrorAndExit(errors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
		}
	}
}

func getRingName() string {

	return utils.GetStoredRingName()
}
