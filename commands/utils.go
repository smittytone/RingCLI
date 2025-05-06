package commands

import (
	"github.com/spf13/cobra"

	utilsCommands "ringcli/commands/utils"
)

// Define the `utils` command.
var utilsCommand = &cobra.Command{
	Use:   "utils",
	Short: "Ring utility commands",
	Long:  "Scan for Colmi R02 rings and perform housekeeping on one of them.",
	Run:   showAppHelp,
	ValidArgs: []string{
		"battery",
		"bind",
		"find",
		"heartrate",
		"info",
		"scan",
		"shutdown",
		"time",
	},
}

func init() {

	utilsCommand.AddCommand(utilsCommands.BatteryCmd)
	utilsCommand.AddCommand(utilsCommands.BindCmd)
	utilsCommand.AddCommand(utilsCommands.FindCmd)
	utilsCommand.AddCommand(utilsCommands.HeartRateCmd)
	utilsCommand.AddCommand(utilsCommands.InfoCmd)
	utilsCommand.AddCommand(utilsCommands.ScanCmd)
	utilsCommand.AddCommand(utilsCommands.ShutdownCmd)
	utilsCommand.AddCommand(utilsCommands.TimeCmd)
	rootCmd.AddCommand(utilsCommand)
}
