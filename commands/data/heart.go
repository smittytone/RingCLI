package rcDataCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
	"time"
	"tinygo.org/x/bluetooth"
)

var (
	heartRateData *ring.HeartRateData // Pointer to received heart rate log data
	dataEnabled   bool                // Heart rate data monitoring is enabled
	dataInterval  int                 // Heart rate data monitoring time interval eg. 30 minutes
)

// Define the `heartrate` sub-command.
var HeartRateCmd = &cobra.Command{
	Use:   "heartrate",
	Short: "Get your current heart rate",
	Long:  "Get your current heart rate",
	Run:   getHeartRate,
}

func getHeartRate(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	log.Prompt("Retrieving heart rate data")

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)

	// Get the data interval -- we'll use this to parse the received data
	ble.RequestDataViaCommandUART(device, ring.MakeHeartRatePeriodGetRequest(), receiveHeartDataSettings, 1)

	// Get the activity data
	requestHeartData(device)

	// Output received ring data
	outputHeartData()
}

func receiveHeartDataSettings(receivedData []byte) {

	if receivedData[0] == ring.COMMAND_HEART_RATE_PERIOD {
		// Extract and type the received data
		var b byte
		dataEnabled, b = ring.ParseHeartRatePeriodResponse(receivedData)
		dataInterval = int(b)

		// Signal data received
		ble.UARTInfoReceived = true
	}
}

func requestHeartData(device bluetooth.Device) {

	// TODO Allow date offset to be added via CLI option
	requestPacket := ring.MakeHeartRateReadRequest(utils.StartToday(time.Now().UTC()))
	ble.RequestDataViaCommandUART(device, requestPacket, receiveHeartData, 1)
}

func receiveHeartData(receivedData []byte) {

	data := ring.ParseHeartRateDataResponse(receivedData, dataInterval)
	if data != nil {
		// Got data
		heartRateData = data

		// Signal data received
		ble.UARTInfoReceived = true
	}
}

func outputHeartData() {

	log.ClearPrompt()

	if heartRateData == nil {
		log.ReportError("No heart rate data received")
		return
	}

	log.Report("Heart Data commencing at %s", heartRateData.Time.String())
	noDataMessageStart, noDataMessageEnd := "", ""
	for _, hrdp := range heartRateData.Rates {
		if hrdp.Time.Before(time.Now()) || hrdp.Time.Equal(time.Now()) {
			if hrdp.Rate == 0 {
				if noDataMessageStart == "" {
					noDataMessageStart = hrdp.Timestamp
				}

				noDataMessageEnd = hrdp.Timestamp
			} else {
				if noDataMessageStart != "" {
					if noDataMessageStart != noDataMessageEnd {
						log.Report("  Ring not worn (or no data available) from %s to %s (UTC)", noDataMessageStart, noDataMessageEnd)
					} else {
						log.Report("  Ring not worn (or no data available) at %s (UTC)", noDataMessageStart)
					}
					noDataMessageStart = ""
					noDataMessageEnd = ""
				}

				log.Report("  %d bpm at %s (UTC)", hrdp.Rate, hrdp.Timestamp)
			}
		}
	}

	if noDataMessageStart != "" {
		log.Report("%s %s (UTC)", noDataMessageStart, noDataMessageEnd)
	}
}
