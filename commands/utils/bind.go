package rcUtilsCommands

import (
	// External code
	"github.com/spf13/cobra"
	// Library code
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcUtils "ringcli/lib/utils"
)

// Globals relevant only to this command
var (
	doOverwrite bool = false
	doShow bool = false
)

// Define the `find` subcommand.
var BindCmd = &cobra.Command{
	Use:   "bind",
	Short: "Store a ring address",
	Long:  "Store a ring's address.",
	Run:   bindRing,
}

func bindRing(cmd *cobra.Command, args []string) {

	if doShow {
		// Just show binding info
		ringAddress = rcUtils.GetStoredRingAddress()
		if ringAddress != "" {
			rcLog.Report("Ring %s is currently bound", ringAddress)
		} else {
			rcLog.Report("No ring bound")
		}
	} else {
		// Set binding
		if ringAddress == "" {
			// Bail when no ring address has been provided
			rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No address supplied")
		}

		// Write out the binding for future use
		rcUtils.MakeBinding(ringAddress, doOverwrite)
		rcLog.Report("Ring %s bound", ringAddress)
	}
}
