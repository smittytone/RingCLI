package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	ring "ringcli/lib/colmi"
	errors "ringcli/lib/errors"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
)

// Globals relevant only to this command
var (
	heartRateEnableSet  bool = false // User has asked to enable heart rate monitoring
	heartRateDisableSet bool = false // User has asked to disable heart rate monitoring
	heartRatePeriod     int  = 60    // Heart rate monitoring period
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

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	// Check params: period in minutes
	if heartRatePeriod < 0 || heartRatePeriod > 255 {
		log.ReportErrorAndExit(errors.ERROR_CODE_BAD_PARAMS, "Heart rate reading period out of range (0-255 minutes)")
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

	bspCount = log.Raw("Setting heart rate monitoring state...  ")
	utils.AnimateCursor()

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)
	ble.RequestDataViaCommandUART(device, ring.MakeHeartRatePeriodSetRequest(enabled, byte(heartRatePeriod)), receiveHeartRatePeriod, 1)
}

func getHeartRatePeriod(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	getRingAddress()

	bspCount = log.Raw("Getting heart rate monitoring state...  ")
	utils.AnimateCursor()

	// Enable BLE
	device := ble.EnableAndConnect(ringAddress)
	defer ble.Disconnect(device)
	ble.RequestDataViaCommandUART(device, ring.MakeHeartRatePeriodGetRequest(), receiveHeartRatePeriod, 1)
}

func receiveHeartRatePeriod(receivedData []byte) {

	if receivedData[0] == ring.COMMAND_HEART_RATE_PERIOD {
		utils.StopAnimation()

		// Parse and report received data
		enabled, period := ring.ParseHeartRatePeriodResponse(receivedData)
		enabledString := "enabled"
		if !enabled {
			enabledString = "disabled"
		}

		log.Backspace(bspCount)
		log.Report("Periodic heart rate monitoring is %s", enabledString)

		// Only output period if periodic readings are enabled and we're making a GET request
		// NOTE Interval not included on a SET request for some reason
		if enabled && receivedData[1] == 1 {
			log.Report("Readings taken every %d minutes", period)
		}

		// Signal data received
		ble.UARTInfoReceived = true
	}
}
