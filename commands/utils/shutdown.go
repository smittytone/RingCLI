package rcUtilsCommands

import (
	"fmt"
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
)

// Define the `shutdown` subcommand.
var ShutdownCmd = &cobra.Command{
	Use:   "shutdown",
	Short: "Power down a ring",
	Long:  "Power down a ring. Connect the ring to its charger to restart it.",
	Run:   shutdownRing,
}

func shutdownRing(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)
	ble.RequestDataViaCommandUART(device, ring.MakeShutdownRequest(), shutdownPacketSent, 0)
}

func shutdownPacketSent(receivedData []byte) {

	// NOTE Will not be called -- ie not transmitted just before shutdown (as you might expect)
	fmt.Println(receivedData)
}
