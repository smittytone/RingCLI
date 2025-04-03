package rcDataCommands

import (
	"fmt"
	"github.com/spf13/cobra"
	rcBLE "ringcli/lib/ble"
	rcColmi "ringcli/lib/colmi"
	rcLog "ringcli/lib/log"
	rcUtils "ringcli/lib/utils"
	"time"
	"tinygo.org/x/bluetooth"
)

var (
	heartRateData *rcColmi.HeartRateData
	dataEnabled   bool
	dataInterval  int
)

// Define the `heartrate` sub-command.
var HeartRateCmd = &cobra.Command{
	Use:   "heartrate",
	Short: "Get your current heart rate",
	Long:  "Get your current heart rate",
	Run:   getHeartRate,
}

func getHeartRate(cmd *cobra.Command, args []string) {

	getRingAddress()

	// Enable BLE
	bspCount = rcLog.Raw("Retrieving data...")
	device := rcBLE.EnableAndConnect(ringAddress)
	defer rcBLE.Disconnect(device)

	// Get the data interval
	rcBLE.RequestDataViaCommandUART(device, rcColmi.MakeHeartRatePeriodGetReq(), receiveHeartDataSettings, 1)

	// Get the activity data
	requestHeartData(device)

	// Output received ring data
	outputHeartData()
}

func requestHeartData(ble bluetooth.Device) {

	// TODO Allow date offset to be added via cli option
	requestPacket := rcColmi.MakeHeartRateReadReq(rcUtils.StartToday(time.Now().UTC()))
	rcBLE.RequestDataViaCommandUART(ble, requestPacket, receiveHeartData, 1)
}

func receiveHeartData(receivedData []byte) {

	//fmt.Println(receivedData)
	a := rcColmi.ParseHeartRateDataResp(receivedData, dataInterval)
	if a != nil {
		// Got data
		rcBLE.UARTInfoReceived = true
		heartRateData = a
	}
}

func receiveHeartDataSettings(receivedData []byte) {

	if receivedData[0] == rcColmi.COMMAND_HEART_RATE_PERIOD {
		// Signal data received
		var b byte
		dataEnabled, b = rcColmi.ParseHeartRatePeriodResp(receivedData)
		dataInterval = int(b)
		rcBLE.UARTInfoReceived = true
	}
}

func outputHeartData() {

	start, end := "", ""
	rcLog.Backspaces(bspCount)
	rcLog.Report("Heart Data commencing at %s", heartRateData.Timestamp.String())
	for _, hrdp := range heartRateData.Rates {
		if hrdp.Timestamp.Before(time.Now()) || hrdp.Timestamp.Equal(time.Now()) {
			if hrdp.Rate == 0 {
				if start == "" {
					start = fmt.Sprintf("  Ring not worn or no data available from %s to", hrdp.Time)
				} else {
					end = hrdp.Time
				}
			} else {
				if start != "" {
					rcLog.Report("%s %s (UTC)", start, end)
					start = ""
					end = "now"
				}

				rcLog.Report("  %d bpm at %s (UTC)", hrdp.Rate, hrdp.Time)
			}
		}
	}

	if start != "" {
		rcLog.Report("%s %s (UTC)", start, end)
	}
}
