package rcUtilsCommands

import (
	"fmt"
	"time"
	// External code
	"github.com/spf13/cobra"
	// Library code
	rcBLE "ringcli/lib/ble"
	rcColmi "ringcli/lib/colmi"
)

// Globals relevant only to this command
var (
	flashCount      uint = 1
	doneFlag        bool = false
	continuousFlash bool = false
)

// Define the `find` subcommand.
var FindCmd = &cobra.Command{
	Use:   "find",
	Short: "Locate ring",
	Long:  "Locate a ring by flashing its green LED",
	Run:   findRing,
}

func findRing(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	fmt.Println("Look for your ring")

	// Set a long flash count for 'continuous', ie. not literally so
	if continuousFlash {
		flashCount = 99
	}

	// Enable BLE
	device := rcBLE.EnableAndConnect(ringAddress)
	defer rcBLE.Disconnect(device)
	requestPacket := rcColmi.MakeLedFlashReq()
	rcBLE.RequestDataViaCommandUART(device, requestPacket, flashPacketResponseReceived, flashCount)
}

func flashPacketResponseReceived(receivedData []byte) {

	if receivedData[0] == rcColmi.COMMAND_BATTERY_FLASH_LED {
		if continuousFlash {
			// Pause between flashes to ensure smooth operation
			time.Sleep(2 * time.Second)
		}

		rcBLE.UARTInfoReceived = true
	}
}
