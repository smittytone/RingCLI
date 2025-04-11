package rcDataCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
	"strings"
	"time"
)

var (
	bloodOxygenData *ring.BloodOxygenData // Pointer to received heart rate log data
)

// Define the `spo2` sub-command.
var BloodOxygenCmd = &cobra.Command{
	Use:   "spo2",
	Short: "Get your current blood oxygen readings",
	Long:  "Get your current blood oxygen readings",
	Run:   getBloodOxygen,
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

	if bloodOxygenData == nil {
		log.ReportError("No blood oxygen data received")
		return
	}

	if showFull {
		log.Report("Full blood oxygen data from %s", bloodOxygenData.Time.String())
		for _, dataSet := range bloodOxygenData.Data {
			log.Report("  %s", dataSet.Time.String())
			outputDataSet(dataSet, "    ")
		}
	} else {
		count := len(bloodOxygenData.Data)
		displaySet := bloodOxygenData.Data[count-1]
		y, m, d := displaySet.Time.Date()
		log.Report("Blood oxygen data for %d-%02d-%02d:", y, m, d)
		outputDataSet(displaySet, "  ")
	}
}

func outputDataSet(dataSet ring.BloodOxygenDataSet, padding string) {

	for _, dataPoint := range dataSet.Rates {
		if dataPoint.Time.Before(time.Now().UTC()) || dataPoint.Time.Equal(time.Now().UTC()) {
			if dataPoint.Maximum < 80 {
				log.Report(padding+"%s (UTC) Ring not worn, or no data available", dataPoint.Timestamp)
			} else {
				log.Report(padding+"%s (UTC) %d%%", dataPoint.Timestamp, dataPoint.Maximum)
			}

			if showFull && dataPoint.Timestamp == "23:00" {
				log.Report(padding + strings.Repeat("-", 47))
			}
		}
	}
}
