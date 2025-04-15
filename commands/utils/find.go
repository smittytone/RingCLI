package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
	"time"
)

// Globals relevant only to this command
var (
	flashCount      uint = 1     // Number of times to flash the ring LED
	doneFlag        bool = false // Operation completed
	continuousFlash bool = false // Flash the LED multiple times
)

// Define the `find` subcommand.
var FindCmd = &cobra.Command{
	Use:   "find",
	Short: "Locate ring",
	Long:  "Help locate a smart ring by flashing its green LED.",
	Run:   findRing,
}

func findRing(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	// Set a long flash count for 'continuous', ie. not literally so
	if continuousFlash {
		flashCount = 99
		log.Report("Ring LED flashing continuously")
	} else {
		log.Report("Ring LED flashing")
	}

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)
	requestPacket := ring.MakeLedFlashRequest()
	ble.RequestDataViaCommandUART(device, requestPacket, flashPacketResponseReceived, flashCount)
}

func flashPacketResponseReceived(receivedData []byte) {

	if receivedData[0] == ring.COMMAND_BATTERY_FLASH_LED {
		if continuousFlash {
			// Pause between flashes to ensure smooth operation
			time.Sleep(2 * time.Second)
		}

		// Mark the packet as received
		// On continuous mode this triggers a resend of the flash command
		ble.UARTInfoReceived = true
	}
}
