package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	errors "ringcli/lib/errors"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
)

// Globals relevant only to this command
var (
	doOverwrite bool = false // Overwrite an existing address
	doShow      bool = false // Display the binding, if present
)

// Define the `bind` subcommand.
var BindCmd = &cobra.Command{
	Use:   "bind",
	Short: "Store a ring BLE address",
	Long:  "Persist your ring's BLE address across commands.",
	Run:   bindRing,
}

func bindRing(cmd *cobra.Command, args []string) {

	if doShow {
		// Just show binding info
		ringAddress = utils.GetStoredRingAddress()
		if ringAddress != "" {
			log.Report("Ring %s is currently bound", ringAddress)
		} else {
			log.Report("No ring bound")
		}
	} else {
		// Set binding
		if ringAddress == "" {
			// Bail when no ring address has been provided
			log.ReportErrorAndExit(errors.ERROR_CODE_BAD_PARAMS, "No address supplied")
		}

		// Write out the binding for future use
		utils.MakeBinding(ringAddress, doOverwrite)
		log.Report("Ring %s bound", ringAddress)
	}
}
