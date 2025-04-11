package rcDataCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
)

var (
	sleepData *ring.SleepData // Pointer to received sleep log data
)

// Define the `sleep` sub-command.
var SleepCmd = &cobra.Command{
	Use:   "sleep",
	Short: "Get your current sleep report",
	Long:  "Retrieve your current sleep report.",
	Run:   getSleepReport,
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

	if sleepData == nil {
		log.ReportError("No sleep data received")
		return
	}

	if showFull {
		for _, period := range sleepData.Periods {
			outputSleepPeriod(period)
		}
	} else {
		count := len(sleepData.Periods)
		displayPeriod := sleepData.Periods[count-1]
		outputSleepPeriod(displayPeriod)
	}
}

func outputSleepPeriod(period ring.SleepPeriod) {

	// Get total sleep duration
	total := 0
	for _, phase := range period.Phases {
		total += phase.Duration
	}

	// Convert minutes duration to hours and minutes, and output
	hours := total / 60
	mins := total - (hours * 60)
	log.Report("Sleep data from %s to %s ", period.StartTime.String(), period.EndTime.String())
	log.Report("  Sleep duration: %d hours, %d minutes comprising:", hours, mins)

	// Output the phase data
	for _, phase := range period.Phases {
		log.Report("    %d minutes %s", phase.Duration, ring.GetSleepType(phase.Type))
	}
}
