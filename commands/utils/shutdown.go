package rcUtilsCommands

import (
	"fmt"
	// External code
	"github.com/spf13/cobra"
	// Library code
	rcBLE "ringcli/lib/ble"
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcColmi "ringcli/lib/colmi"
)

// Define the `scan` subcommand.
var ShutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Power down a ring",
	Long:  "Power down a ring. Connect the ring to its charge to restart it.",
	Run:   shutdownRing,
}

func shutdownRing(cmd *cobra.Command, args []string) {

	// Bail when no ID data is provided
	if ringName == "" && ringAddress == "" {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
	}

	// Enable BLE
	device := rcBLE.Connect(rcBLE.Open(), rcBLE.AddressFromString(ringAddress))
	defer rcBLE.Disconnect(device)
	rcBLE.RequestDataViaCommandUART(device, rcColmi.MakeShutdownReq(), shutdownPacketSent, 0)
}

func shutdownPacketSent(receivedData []byte) {

	// NOTE Will not be called
	fmt.Println(receivedData)
}