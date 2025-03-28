package rcDataCommands

import (
	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"

	rcBLE "ringcli/lib/ble"
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcColmi "ringcli/lib/colmi"
)

var (
	activityInfo rcColmi.SportsInfo
)

// Define the `scan` subcommand.
var StepsCmd = &cobra.Command{
	Use:   "steps",
	Short: "Get activity info",
	Long: "Get activity info",
	Run:    getSteps,
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

	// Get the activity data
	requestSportsInfo(device)

	// Output received ring data
	outputInfo()
}

func requestSportsInfo(ble bluetooth.Device) {

	// TODO Allow date offset to be added via cli option
	requestPacket := rcColmi.MakeStepsReq(0)
	if requestPacket == nil {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_ACTIVITY_TIME_OFFSET, "Date offset out of range")
	}

	activityInfo = rcColmi.SportsInfo{}
	rcBLE.RequestDataViaUART(ble, requestPacket, receiveSportsInfo)
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
		activityInfo.Timestamp.Year = info.Timestamp.Year
		activityInfo.Timestamp.Month = info.Timestamp.Month
		activityInfo.Timestamp.Day = info.Timestamp.Day
		activityInfo.Timestamp.Hour = info.Timestamp.Hour
		activityInfo.Timestamp.Minutes = info.Timestamp.Minutes
		activityInfo.NewCalories = info.NewCalories
		activityInfo.Calories += info.Calories
		activityInfo.Steps += info.Steps
		activityInfo.Distance += info.Distance

		//now := time.Now().Hour() * 4
		if info.IsDone {
			// Completion notice -- but we'll probably have timed out before this
			rcBLE.UARTInfoReceived = true
		}
	}
}

func outputInfo() {

	rcLog.Report("Activity Info: %d/%d/%d %02d:%02d", activityInfo.Timestamp.Day, activityInfo.Timestamp.Month, activityInfo.Timestamp.Year, activityInfo.Timestamp.Hour, activityInfo.Timestamp.Minutes)
	rcLog.Report("         Steps: %d", activityInfo.Steps)

	// Check for later, alternative calories scaling
	if activityInfo.NewCalories {
		rcLog.Report("      Calories: %.02f kCal", float32(activityInfo.Calories) / 1000)
	} else {
		rcLog.Report("      Calories: %d kCal", activityInfo.Calories)
	}

	// Adjust for range of movement order of magnitude
	if activityInfo.Distance > 999 {
		rcLog.Report("Distance Moved: %.02f km", float32(activityInfo.Distance) / 1000)
	} else {
		rcLog.Report("Distance Moved: %d m", activityInfo.Distance)
	}
}