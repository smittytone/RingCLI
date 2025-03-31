package rcUtilsCommands

import (
	// External code
	"github.com/spf13/cobra"
	// Library code
	rcBLE "ringcli/lib/ble"
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
	rcColmi "ringcli/lib/colmi"
)

// Globals relevant only to this command
var (
	heartRateEnableSet 	bool = false
	heartRateDisableSet bool = false
	heartRatePeriod     int = 60
)

// Define the `setheartrate` sub-command.
var SetHeartRateCmd = &cobra.Command{
	Use:   "setheartrate",
	Short: "Set the heart rate monitoring period",
	Long:  "Set the heart rate monitoring period in minutes, and enable or disable monitoring.",
	Run:   setHeartRatePeriod,
}

// Define the `getheartrate` sub-command.
var GetHeartRateCmd = &cobra.Command{
	Use:   "getheartrate",
	Short: "Get the heart rate monitoring period",
	Long:  "Get the heart rate monitoring state and, if enabled, its periodicity in minutes.",
	Run:   getHeartRatePeriod,
}

func setHeartRatePeriod(cmd *cobra.Command, args []string) {

	// Bail when no ID data is provided
	if ringName == "" && ringAddress == "" {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
	}

	// Check params: period in minutes
	if heartRatePeriod < 0 || heartRatePeriod > 255 {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "Heart rate reading period out of range (0-255 minutes)")
	}

	// Check params: enable periodic readings, default to `true`
	enabled := true
	if heartRateDisableSet {
		// `--disable` specified
		enabled = false
	}

	if heartRateEnableSet {
		// `--enable` specified and overrides `--disable` if also set
		enabled = true
	}

	if heartRatePeriod == 0 {
		// Setting the period to zero overrides `--enable`
		enabled = false
		heartRatePeriod = 60
	}

	// Enable BLE
	device := rcBLE.EnableAndConnect(ringAddress)
	defer rcBLE.Disconnect(device)
	rcBLE.RequestDataViaCommandUART(device, rcColmi.MakeHeartRatePeriodSetReq(enabled, byte(heartRatePeriod)), heartRatePeriodPacketReceived, 1)
}

func getHeartRatePeriod(cmd *cobra.Command, args []string) {

	// Bail when no ID data is provided
	if ringName == "" && ringAddress == "" {
		rcLog.ReportErrorAndExit(rcErrors.ERROR_CODE_BAD_PARAMS, "No name or address supplied")
	}

	// Enable BLE
	device := rcBLE.EnableAndConnect(ringAddress)
	defer rcBLE.Disconnect(device)
	rcBLE.RequestDataViaCommandUART(device, rcColmi.MakeHeartRatePeriodGetReq(), heartRatePeriodPacketReceived, 1)
}

func heartRatePeriodPacketReceived(receivedData []byte) {

	if receivedData[0] == rcColmi.COMMAND_HEART_RATE_PERIOD {
		// Signal data received
		rcBLE.UARTInfoReceived = true

		// Parse and report received data
		enabled, period := rcColmi.ParseHeartRatePeriodResp(receivedData)
		enabledString := "enabled"
		if !enabled {
			enabledString = "disabled"
		}

		rcLog.Report("Periodic heart rate monitoring is %s", enabledString)

		// Only output period if periodic readings are enabled and we're making a GET request
		// NOTE Interval not included on a SET request for some reason
		if enabled && receivedData[1] == 1 {
			rcLog.Report("Readings taken every %d minutes", period)
		}
	}
}