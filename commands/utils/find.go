package rcUtilsCommands

import (
	"fmt"
	"time"
	// External code
	"github.com/spf13/cobra"
	// Library code
	rcBLE "ringcli/lib/ble"
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcColmi "ringcli/lib/colmi"
)

var (
	doneFlag bool = false
	flashCount int = 1
)

// Define the `scan` subcommand.
var FindCmd = &cobra.Command{
	Use:   "find",
	Short: "Locate ring",
	Long:  "Locate a ring by flashing its green LED",
	Run:   findRing,
}

func findRing(cmd *cobra.Command, args []string) {

	// Bail when no ID data is provided
	if ringName == "" && ringAddress == "" {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
	}

	fmt.Println("Look for your ring")

	if continuousFlash {
		flashCount = 99
	}

	// Enable BLE
	ble := rcBLE.Open()
	bleAddress := rcBLE.AddressFromString(ringAddress)
	device := rcBLE.Connect(ble, bleAddress)
	defer rcBLE.Disconnect(device)
	requestPacket := rcColmi.MakeLedFlashReq()
	rcBLE.RequestDataViaCommandUART(device, requestPacket, packetSent, flashCount)
}

func packetSent(receivedData []byte) {

	if receivedData[0] == rcColmi.COMMAND_BATTERY_FLASH_LED {
		if continuousFlash {
			// Pause between flashes to ensure smooth operation
			time.Sleep(2 * time.Second)
		}
		rcBLE.UARTInfoReceived = true
	}
}