package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
	"time"
)

// Define the `settime` subcommand.
var SetTimeCmd = &cobra.Command{
	Use:   "time",
	Short: "Initialise the ring's date and time",
	Long:  "Initialise the ring's date and time.",
	Run:   setTime,
}

func setTime(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	log.Prompt("Setting your ring's date and time")

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)
	requestPacket := ring.MakeTimeSetRequest(time.Now())
	ble.RequestDataViaCommandUART(device, requestPacket, setTimeResponseReceived, 1)

	log.ClearPrompt()
	log.Report("Ring's time set")
}

func setTimeResponseReceived(receivedData []byte) {

	if receivedData[0] == ring.COMMAND_SET_TIME {
		ble.UARTInfoReceived = true
	}
}
