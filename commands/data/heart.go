package rcDataCommands

import (
	"fmt"
	"time"
	// External code
	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"
	// Library code
	rcBLE "ringcli/lib/ble"
	rcColmi "ringcli/lib/colmi"
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcUtils "ringcli/lib/utils"
)

var (
	heartRateData *rcColmi.HeartRateData
)

// Define the `heartrate` sub-command.
var HeartRateCmd = &cobra.Command{
	Use:   "heartrate",
	Short: "Get your current heart rate",
	Long:  "Get your current heart rate",
	Run:   getHeartRate,
}

func getHeartRate(cmd *cobra.Command, args []string) {

	// Bail when no ID data is provided
	if ringAddress == "" {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
	}

	// Enable BLE
	bspCount = rcLog.Raw("Retrieving data...")
	device := rcBLE.EnableAndConnect(ringAddress)
	defer rcBLE.Disconnect(device)

	// Get the activity data
	requestHeartData(device)

	// Output received ring data
	outputHeartData()
}

func requestHeartData(ble bluetooth.Device) {

	// TODO Allow date offset to be added via cli option
	requestPacket := rcColmi.MakeHeartRateReadReq(rcUtils.StartToday(time.Now().UTC()))
	fmt.Println(requestPacket)
	rcBLE.RequestDataViaCommandUART(ble, requestPacket, receiveHeartData, 1)
}

func receiveHeartData(receivedData []byte) {

	fmt.Println(receivedData)
	a := rcColmi.ParseHeartRateDataResp(receivedData)
	if a != nil {
		// Got data
		rcBLE.UARTInfoReceived = true
		heartRateData = a
	} else {
		fmt.Println("NIL")
	}
}

func outputHeartData() {

	rcLog.Report("Heart Data from %s (%d)", heartRateData.Timestamp.String(), len(heartRateData.Rates))
	for _, hrdp := range heartRateData.Rates {
		rcLog.Report("  %d bpm at %s", hrdp.Rate, hrdp.Timestamp)
	}
}
