package rcDataCommands

import (
	"fmt"
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
	"time"
)

var (
	heartRateData         *ring.HeartRateData // Pointer to received heart rate log data
	dataEnabled           bool                // Heart rate data monitoring is enabled
	dataInterval          int                 // Heart rate data monitoring time interval eg. 30 minutes
	heartRateDataRealtime []int               // Real time value store
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

	// Set the terminal text
	if inRealTime {
		log.Prompt("Setting up real-time heart rate monitoring (can take up to 30 seconds)")
	} else {
		log.Prompt("Retrieving heart rate data")
	}

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)

	if inRealTime {
		// Poll for heart rate data in real time: readings every 30s
		ble.PollRealtime(device, ring.REAL_TIME_HEART_RATE_BATCH, receiveHeartDataRealtime, 10)

		// Output received ring data
		outputHeartDataRealtime()
	} else {
		// Get the data interval -- we'll use this to parse the received data
		ble.RequestDataViaCommandUART(device, ring.MakeHeartRatePeriodGetRequest(), receiveHeartDataSettings, 1)

		// Get the activity data
		ble.RequestDataViaCommandUART(device, ring.MakeHeartRateReadRequest(utils.StartToday(time.Now().UTC())), receiveHeartDataLog, 1)

		// Output received ring data
		outputHeartDataLog()
	}
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

func receiveHeartDataLog(receivedData []byte) {

	if receivedData[0] == ring.COMMAND_HEART_RATE_READ {
		data := ring.ParseHeartRateDataResponse(receivedData, dataInterval)
		if data != nil {
			// Got data
			heartRateData = data

			// Signal data received
			ble.UARTInfoReceived = true
		}
	}
}

func receiveHeartDataRealtime(receivedData []byte) {

	//fmt.Println(receivedData)
	if receivedData[0] == ring.COMMAND_START_REAL_TIME {
		ok, data := ring.ParseRealtimeHeartDataResponse(receivedData)
		if ok {
			log.ClearPrompt()

			// Got data
			if heartRateDataRealtime == nil {
				heartRateDataRealtime = make([]int, 0, 60)
			}

			// Output realtime reading
			heartRateDataRealtime = append(heartRateDataRealtime, data.Value)
			var formatString string
			if ble.PollCount%2 == 0 {
				formatString = "%d bpm ❤️"
			} else {
				formatString = "%d bpm   "
			}

			log.RealtimeDataOut(fmt.Sprintf(formatString, data.Value))

			// Signal data received
			ble.PollCount += 1
		}
	}
}

func outputHeartDataLog() {

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

func outputHeartDataRealtime() {

	log.RealtimeDataClear()

	count := len(heartRateDataRealtime)
	if count > 0 {
		log.Report("Last reading: %d bpm", heartRateDataRealtime[count-1])
		log.Report("Average over %d readings: %d bpm", count, getAverage())
	}
}

func getAverage() int {

	total := 0
	for _, i := range heartRateDataRealtime {
		total += i
	}

	return total / len(heartRateDataRealtime)
}
