package rcUtilsCommands

import (
	"fmt"
	// External code
	"github.com/spf13/cobra"
	// Library code
	rcBLE "ringcli/lib/ble"
	//rcErrors "ringcli/lib/errors"
	//rcLog "ringcli/lib/log"
	rcColmi "ringcli/lib/colmi"
)

// Define the `shutdown` subcommand.
var ShutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Power down a ring",
	Long:  "Power down a ring. Connect the ring to its charge to restart it.",
	Run:   shutdownRing,
}

func shutdownRing(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	// Enable BLE
	device := rcBLE.EnableAndConnect(ringAddress)
	defer rcBLE.Disconnect(device)
	rcBLE.RequestDataViaCommandUART(device, rcColmi.MakeShutdownReq(), shutdownPacketSent, 0)
}

func shutdownPacketSent(receivedData []byte) {

	// NOTE Will not be called
	fmt.Println(receivedData)
}