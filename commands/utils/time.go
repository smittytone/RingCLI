package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	rcBLE "ringcli/lib/ble"
	rcColmi "ringcli/lib/colmi"
	rcLog "ringcli/lib/log"
	rcUtils "ringcli/lib/utils"
	"time"
)

// Globals relevant only to this command
var (
	timeSet bool = false
)

// Define the `set time` subcommand.
var SetTimeCmd = &cobra.Command{
	Use:   "settime",
	Short: "Initialise the ring's date and time",
	Long:  "Initialise the ring's date and time.",
	Run:   setTime,
}

func setTime(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	rcLog.Report("Setting your ring's date and time...  ")
	rcUtils.AnimateCursor()

	// Enable BLE
	device := rcBLE.EnableAndConnect(ringAddress)
	defer rcBLE.Disconnect(device)
	requestPacket := rcColmi.MakeTimeSetReq(time.Now())
	rcBLE.RequestDataViaCommandUART(device, requestPacket, setTimeResponseReceived, 1)

	if timeSet {
		rcLog.Report("Ring's time set")
	}
}

func setTimeResponseReceived(receivedData []byte) {

	if receivedData[0] == rcColmi.COMMAND_SET_TIME {
		rcUtils.StopAnimation()
		timeSet = true
		rcBLE.UARTInfoReceived = true
	}
}
