package rcDataCommands

import (
	// External code
	"github.com/spf13/cobra"
	// App
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcUtils "ringcli/lib/utils"
)

// Globals for all commands
var (
	debug       bool = false
	ringName    string = ""
	ringAddress string = ""
	bspCount    int = 0
)

// Set up the `data` sub-commands' flags.
func init() {
	// Add optional flags: --address
	StepsCmd.Flags().StringVarP(&ringAddress, "address", "", "", "The rings's BLE address. Required")

	// Add required flags: --address
	HeartRateCmd.Flags().StringVarP(&ringAddress, "address", "", "", "The rings's BLE address. Required")
	HeartRateCmd.MarkFlagRequired("address")
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
		ringAddress = rcUtils.GetStoredRingAddress()

		if ringAddress == "" {
			// No loaded address so report and bail
			rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
		}
	}
}