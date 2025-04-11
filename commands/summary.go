package commands

import (
	"github.com/spf13/cobra"
	ble "ringcli/lib/ble"
	cmds "ringcli/commands/utils"
	data "ringcli/commands/data"
	log "ringcli/lib/log"
	"strings"
)

// Globals
var (
	debug       bool   = false // Are we in debug mode?
	//RingAddress string = ""    // Ring BLE address
	deviceInfo cmds.DeviceInfo = cmds.DeviceInfo{}
)

// Define the `utils` command.
var rcSummaryCommand = &cobra.Command{
	Use:   "summary",
	Short: "Get a daily data summary",
	Long:  "Present a daily data summary.",
	Run:   showSummary,
}

func init() {

	rootCmd.AddCommand(rcSummaryCommand)
}

func showSummary(cmd *cobra.Command, args []string) {

	// Make sure we have a ring BLE address from the command line or store
	cmds.GetRingAddress()

	log.Prompt("Preparing summary")

	// Enable BLE
	device := ble.EnableAndConnect(cmds.RingAddress)
	defer ble.Disconnect(device)
	cmds.RequestBatteryInfo(device)
	data.RequestSportsInfo(device)
	data.RequestHeartData(device)
	data.RequestBloodOxygenData(device)

	log.Report(strings.Repeat("*", 10) + " DEVICE " + strings.Repeat("*", 10))
	cmds.OutputBatteryInfo()
	log.Report(strings.Repeat("*", 9) + " ACTIVITY " + strings.Repeat("*", 9))
	data.OutputStepsInfo()
	log.Report(strings.Repeat("*", 10) + " HEALTH " + strings.Repeat("*", 10))
	data.OutputHeartData()
	log.Report(" ")
	data.OutputBloodOxygenData()
}
