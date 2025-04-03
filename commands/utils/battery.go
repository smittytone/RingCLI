package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	rcBLE "ringcli/lib/ble"
	rcLog "ringcli/lib/log"
	rcUtils "ringcli/lib/utils"
)

// Define the `info` subcommand.
var BatteryCmd = &cobra.Command{
	Use:   "battery",
	Short: "Get ring battery state",
	Long:  "Get ring battery state",
	Run:   getBatteryState,
}

func getBatteryState(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	bspCount = rcLog.Raw("Retrieving ring battery state...  ")
	rcUtils.AnimateCursor()

	// Enable BLE
	deviceInfo.battery.Level = 0
	device := rcBLE.EnableAndConnect(ringAddress)
	defer rcBLE.Disconnect(device)
	requestBatteryInfo(device)

	// Output received ring data
	outputRingInfo(true)
}
