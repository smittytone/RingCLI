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
	ringAddress      string   = ""    // Ring BLE address
)

// Set up the `data` sub-commands' flags.
func init() {
	// Add optional flags: --address
	StepsCmd.Flags().StringVarP(&ringAddress, "address", "", "", ADDRESS_COMMAND_TEXT)

	// Add optional flags: --address
	HeartRateCmd.Flags().StringVarP(&ringAddress, "address", "", "", ADDRESS_COMMAND_TEXT)
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
