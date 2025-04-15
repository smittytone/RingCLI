package rcDataCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	errors "ringcli/lib/errors"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
	"tinygo.org/x/bluetooth"
)

// Globals relevant only to this command
var (
	activityTotals ring.ActivityInfo // Values combined from individual records
)

// Define the `steps` sub-command.
var StepsCmd = &cobra.Command{
	Use:   "steps",
	Short: "Get activity info",
	Long:  "Get activity info",
	Run:   getSteps,
}

func getSteps(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	log.Prompt("Retrieving activity data")

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)

	// Get the activity data
	requestActivityInfo(device)

	// Output received ring data
	outputStepsInfo()
}

func requestActivityInfo(device bluetooth.Device) {

	// TODO Allow date offset to be added via cli option
	requestPacket := ring.MakeStepsRequest(0)
	if requestPacket == nil {
		log.ReportErrorAndExit(errors.ERROR_CODE_BAD_ACTIVITY_TIME_OFFSET, "Date offset out of range")
	}

	ts := ring.TimestampFromNow()
	activityTotals = ring.ActivityInfo{
		Steps:       0,
		Distance:    0,
		Calories:    0,
		NoData:      true,
		IsDone:      false,
		NewCalories: false,
		Timestamp:   ts,
	}

	ble.RequestDataViaCommandUART(device, requestPacket, receiveActivityInfo, 1)
}

func receiveActivityInfo(receivedData []byte) {

	if receivedData[0] == ring.COMMAND_GET_ACTIVITY_UNKNOWN {
		// Completion notice???
		ble.UARTInfoReceived = true
	}

	// Check we have a ActivityInfo response
	if receivedData[0] == ring.COMMAND_GET_ACTIVITY_DATA {
		info := ring.ParseStepsResponse(receivedData)
		if info.NoData {
			// Mark as done
			ble.UARTInfoReceived = true
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

		if info.IsDone {
			// Completion notice -- but we'll probably have timed out before this
			ble.UARTInfoReceived = true
		}
	}
}

func outputStepsInfo() {

	log.ClearPrompt()

	// Output...
	log.Report("Activity Info for %d %s %d:", activityTotals.Timestamp.Day, utils.StringifyMonth(activityTotals.Timestamp.Month), activityTotals.Timestamp.Year)
	log.Report("  👟 %d", activityTotals.Steps)

	// Check for later, alternative calories scaling
	if activityTotals.NewCalories {
		log.Report("  ⚡️ %.02f kCal", float32(activityTotals.Calories)/1000)
	} else {
		log.Report("  ⚡️ %d kCal", activityTotals.Calories)
	}

	// Adjust for range of movement order of magnitude
	if activityTotals.Distance > 999 {
		log.Report("  📏 %.02f km", float32(activityTotals.Distance)/1000)
	} else {
		log.Report("  📏 %d m", activityTotals.Distance)
	}
}
