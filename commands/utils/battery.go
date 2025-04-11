package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	log "ringcli/lib/log"
)

// Define the `battery` subcommand.
var BatteryCmd = &cobra.Command{
	Use:   "battery",
	Short: "Get ring battery state",
	Long:  "Retrieve the battery state of your smart ring.",
	Run:   getBatteryState,
}

func getBatteryState(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	log.Prompt("Retrieving ring battery state")

	// Enable BLE
	deviceInfo.battery.Level = 0
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)
	requestBatteryInfo(device)

	// Output received ring data
	outputBatteryInfo()
}

func outputBatteryInfo() {

	chargeState := getChargeState(deviceInfo.battery.IsCharging)
	log.ClearPrompt()
	log.Report("Battery state: %d%% (%s)", deviceInfo.battery.Level, chargeState)
}
