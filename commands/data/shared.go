package rcDataCommands

import (
	"github.com/spf13/cobra"
	errors "ringcli/lib/errors"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
)

const (
	ADDRESS_COMMAND_TEXT = "The rings's BLE address. Required if no ring has been bound."
)

// Globals for all commands
var (
	debug            bool     = false // Are we in debug mode?
	RingAddress      string   = ""    // Ring BLE address
)

// Set up the `data` sub-commands' flags.
func init() {
	// Add optional flags: --address
	StepsCmd.Flags().StringVarP(&RingAddress, "address", "", "", ADDRESS_COMMAND_TEXT)

	// Add optional flags: --address
	HeartRateCmd.Flags().StringVarP(&RingAddress, "address", "", "", ADDRESS_COMMAND_TEXT)

	// Add optional flags: --address, --full
	BloodOxygenCmd.Flags().StringVarP(&RingAddress, "address", "", "", ADDRESS_COMMAND_TEXT)
	BloodOxygenCmd.Flags().BoolVarP(&showFull, "full", "f", false, "Show all available data, not just the most recent 24-hour period.")
}

// Apply the logging Level string, eg. "debug" from the `--log` flag.
// NOTE We have to do this here as it's the first time the flag data
// becomes available to the code.
func processFlags(cmd *cobra.Command, args []string) {

	/*
		err := mvAppConfig.ProcessCommonFlags()
		if err != nil {
			mvLog.ReportErrorAndExit("%s", err)
		}
	*/
}

func GetRingAddress() {

	// Check that `RingAddress` has been set by option: if it has, `RingAddress` will not be empty
	if RingAddress == "" {
		// Try to get a stored (bound)
		RingAddress = utils.GetStoredRingAddress()

		if RingAddress == "" {
			// No loaded address so report and bail
			log.ReportErrorAndExit(errors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
		}
	}
}
