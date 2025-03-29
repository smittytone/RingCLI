package rcDataCommands

import (
	// External code
	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"
	// Library code
	rcBLE "ringcli/lib/ble"
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcColmi "ringcli/lib/colmi"
)

var (
	activityTotals rcColmi.SportsInfo
)

// Define the `scan` subcommand.
var StepsCmd = &cobra.Command{
	Use:   "steps",
	Short: "Get activity info",
	Long:  "Get activity info",
	Run:   getSteps,
}

func getSteps(cmd *cobra.Command, args []string) {

	// Bail when no ID data is provided
	if ringAddress == "" {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
	}

	// Enable BLE
	ble := rcBLE.Open()

	// Generate the ring BLE address and connect to it
	bleAddress := rcBLE.AddressFromString(ringAddress)
	device := rcBLE.Connect(ble, bleAddress)
	defer rcBLE.Disconnect(device)

	// Get the activity data
	requestSportsInfo(device)

	// Output received ring data
	outputStepsInfo()
}

func requestSportsInfo(ble bluetooth.Device) {

	// TODO Allow date offset to be added via cli option
	requestPacket := rcColmi.MakeStepsReq(0)
	if requestPacket == nil {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_ACTIVITY_TIME_OFFSET, "Date offset out of range")
	}

	activityTotals = rcColmi.SportsInfo{}
	rcBLE.RequestDataViaUART(ble, requestPacket, receiveSportsInfo, 1)
}

func receiveSportsInfo(receivedData []byte) {

	if receivedData[0] == 0x73 {
		// Completion notice???
		rcBLE.UARTInfoReceived = true
	}

	// Check we have a SportsInfo response
	if receivedData[0] == 0x43 {
		info := rcColmi.ParseStepsResp(receivedData)
		if info.NoData {
			// Mark as done
			rcBLE.UARTInfoReceived = true
			return
		}

		// Record the data from the packet (one of many)
		activityTotals.Timestamp.Year = info.Timestamp.Year
		activityTotals.Timestamp.Month = info.Timestamp.Month
		activityTotals.Timestamp.Day = info.Timestamp.Day
		activityTotals.Timestamp.Hour = info.Timestamp.Hour
		activityTotals.Timestamp.Minutes = info.Timestamp.Minutes
		activityTotals.NewCalories = info.NewCalories
		activityTotals.Calories += info.Calories
		activityTotals.Steps += info.Steps
		activityTotals.Distance += info.Distance

		//now := time.Now().Hour() * 4
		if info.IsDone {
			// Completion notice -- but we'll probably have timed out before this
			rcBLE.UARTInfoReceived = true
		}
	}
}

func outputStepsInfo() {

	rcLog.Report("Activity Info: %d/%d/%d %02d:%02d", activityTotals.Timestamp.Day, activityTotals.Timestamp.Month, activityTotals.Timestamp.Year, activityTotals.Timestamp.Hour, activityTotals.Timestamp.Minutes)
	rcLog.Report("         Steps: %d", activityTotals.Steps)

	// Check for later, alternative calories scaling
	if activityTotals.NewCalories {
		rcLog.Report("      Calories: %.02f kCal", float32(activityTotals.Calories) / 1000)
	} else {
		rcLog.Report("      Calories: %d kCal", activityTotals.Calories)
	}

	// Adjust for range of movement order of magnitude
	if activityTotals.Distance > 999 {
		rcLog.Report("Distance Moved: %.02f km", float32(activityTotals.Distance) / 1000)
	} else {
		rcLog.Report("Distance Moved: %d m", activityTotals.Distance)
	}
}