package rcDataCommands

import (
	//"fmt"
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
	//"strings"
	//"time"
)

var (
	sleepData *ring.SleepData // Pointer to received sleep log data
)

// Define the `heartrate` sub-command.
var SleepCmd = &cobra.Command{
	Use:       "sleep",
	Short:     "Get your current sleep report",
	Long:      "Get your current sleep report",
	Run:       getSleepReport,
}

func getSleepReport(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	log.Prompt("Retrieving sleep data")

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)

	// Get the data
	ble.RequestDataViaDataUART(device, ring.MakeSleepGetRequest(), receivedSleepData)

	// Output received ring data
	outputSleepData()
}

func receivedSleepData(receivedData []byte) {

	data := ring.ParseSleepDataResponse(receivedData)
	if data != nil {
		// Got data
		sleepData = data

		// Signal data received
		ble.UARTInfoReceived = true
	}
}

func outputSleepData() {

	log.ClearPrompt()

	if showFull {
		log.Report("Full blood oxygen data from %s", bloodOxygenData.Time.String())
	} else {
		count := len(sleepData.Data)
		displayPeriod := sleepData.Data[count - 1]
		total := 0
		for _, period := range displayPeriod.Sleep {
			total += period.Duration
		}

		hours := total / 60
		mins := total - (hours * 60)

		log.Report("Sleep data from %s to %s ", displayPeriod.StartTime.String(), displayPeriod.EndTime.String())
		log.Report("  Sleep duration: %d hours, %d minutes comprising:", hours, mins)

		for _, period := range displayPeriod.Sleep {
			log.Report("    %d minutes %s", period.Duration, ring.GetSleepType(period.Type))
		}


	}
}
