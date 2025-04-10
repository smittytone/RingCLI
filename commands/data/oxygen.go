package rcDataCommands

import (
	"fmt"
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
	"strings"
)

var (
	bloodOxygenData *ring.BloodOxygenData // Pointer to received heart rate log data
)

// Define the `heartrate` sub-command.
var BloodOxygenCmd = &cobra.Command{
	Use:       "spo2",
	Short:     "Get your current blood oxygen readings",
	Long:      "Get your current blood oxygen readings",
	Run:       getBloodOxygen,
}

func getBloodOxygen(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	log.Prompt("Retrieving blood oxygen data")

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)

	// Get the data
	ble.RequestDataViaDataUART(device, ring.MakeBloodOxygenGetRequest(), receiveBloodOxygenData)

	// Output received ring data
	outputBloodOxygenData()
}

func receiveBloodOxygenData(receivedData []byte) {

	fmt.Println(receivedData)
	data := ring.ParseBloodOxygenDataResponse(receivedData)
	if data != nil {
		// Got data
		bloodOxygenData = data

		// Signal data received
		ble.UARTInfoReceived = true
	}
}

func outputBloodOxygenData() {

	log.ClearPrompt()
	log.Report("Blood oxygen data commencing at %s", bloodOxygenData.Time.String())
	for _, bop := range bloodOxygenData.Rates {
		if bop.Maximum < 80 {
			log.Report("  %s (UTC) Ring not worn, or no data available", bop.Timestamp)
		} else {
			log.Report("  %s (UTC) %d%%", bop.Timestamp, bop.Maximum)
		}

		if bop.Timestamp == "23:00" {
			log.Report("  " + strings.Repeat("-", 47))
		}
	}
}
