package rcDataCommands

import (


	"github.com/spf13/cobra"
)

// Globals to hold the device ID etc.
var (
	debug         bool = false
	ringAddress   string = ""
)

// Set up the `utils` command's flags.
func init() {
	// Add required flags: --name
	StepsCmd.Flags().StringVarP(&ringAddress, "address", "", "", "The rings's BLE address. Required")
	StepsCmd.MarkFlagRequired("address")
}

/*
// Check positional args for the commands `info` and `assign`
func validateInfoArgs(cmd *cobra.Command, args []string) error {

	// Add positional args:
	// 1. The device ID
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		// Fails on arg count == 0
		return fmt.Errorf("no name or address provided")
	}

	ringId := args[0]


	if mvHttp.ValidateId(deviceId, utils.DEVICE) == mvSharedData.ID_VALIDATE_ERROR_IS_BAD {
		mvLog.ReportErrorAndExit("Malformed device ID")
	}

	return nil
}
*/

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