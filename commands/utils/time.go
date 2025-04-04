package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
	"time"
)

// Define the `settime` subcommand.
var SetTimeCmd = &cobra.Command{
	Use:   "settime",
	Short: "Initialise the ring's date and time",
	Long:  "Initialise the ring's date and time.",
	Run:   setTime,
}

func setTime(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	bspCount = log.Raw("Setting your ring's date and time...  ")
	utils.AnimateCursor()

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)
	requestPacket := ring.MakeTimeSetRequest(time.Now())
	ble.RequestDataViaCommandUART(device, requestPacket, setTimeResponseReceived, 1)

	log.Backspaces(bspCount)
	log.Report("Ring's time set")
}

func setTimeResponseReceived(receivedData []byte) {

	if receivedData[0] == ring.COMMAND_SET_TIME {
		utils.StopAnimation()
		ble.UARTInfoReceived = true
	}
}
