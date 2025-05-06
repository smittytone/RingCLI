package UtilsCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	config "ringcli/lib/config"
	log "ringcli/lib/log"
)

// Define the `battery` subcommand.
var BatteryCmd = &cobra.Command{
	Use:    "battery",
	Short:  "Get ring battery state",
	Long:   "Retrieve the battery state of your smart ring.",
	PreRun: processGenericFlags,
	Run:    getBatteryState,
}

func getBatteryState(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	log.Prompt("Retrieving ring battery state")

	// Enable BLE
	deviceInfo.battery.Level = 0
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)
	requestBatteryInfo(device) // See `info.go`

	// Output received ring data
	outputBatteryInfo()
}

func outputBatteryInfo() {

	chargeState := getChargeState(deviceInfo.battery.IsCharging) // See `info.go`
	if config.Config.OutputToStdout {
		// Output raw data to stdout
		log.ToStdout("%d", deviceInfo.battery.Level)
	} else if config.Config.OutputToJson {
		// Output data in JSON form to stdout
		log.ToStdout("{\"%s\":{\"battery\":%d}}", ringAddress, deviceInfo.battery.Level)
	} else {
		// Output human readable info to stderr
		log.ClearPrompt()
		log.Report("Battery state: %d%% (%s)", deviceInfo.battery.Level, chargeState)
	}
}
