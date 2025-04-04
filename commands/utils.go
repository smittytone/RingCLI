package commands

import (
	"github.com/spf13/cobra"

	rcUtilsCommands "ringcli/commands/utils"
)

// Define the `utils` command.
var rcUtilsCommand = &cobra.Command{
	Use:   "utils",
	Short: "Ring utility commands",
	Long:  "Scan for Colmi R02 rings and perform housekeeping on one of them.",
	Run:   showAppHelp,
}

func init() {

	rcUtilsCommand.AddCommand(rcUtilsCommands.BatteryCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.BindCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.FindCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.GetHeartRateCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.SetHeartRateCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.InfoCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.ScanCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.ShutdownCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.SetTimeCmd)
	rootCmd.AddCommand(rcUtilsCommand)
}
