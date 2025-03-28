package rcDataCommands

import (
	// External code
	"github.com/spf13/cobra"
)

// Globals for state
var (
	debug       bool = false
	ringName    string = ""
	ringAddress string = ""
)

// Set up the `data` sub-commands' flags.
func init() {
	// Add required flags: --address
	StepsCmd.Flags().StringVarP(&ringAddress, "address", "", "", "The rings's BLE address. Required")
	StepsCmd.MarkFlagRequired("address")
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