package commands

import (
	"github.com/spf13/cobra"

	rcUtilsCommands "ringcli/commands/utils"
)

// Define the `utils` command.
var rcUtilsCommand = &cobra.Command{
	Use:   "utils",
	Short: "Ring utility commands",
	Long:  "Scan for Colmi R02 rings and get current info about one of them.",
	Run:   showAppHelp,
}

func init() {

	rcUtilsCommand.AddCommand(rcUtilsCommands.ScanCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.InfoCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.FindCmd)
	rcUtilsCommand.AddCommand(rcUtilsCommands.ShutdownCmd)
	rootCmd.AddCommand(rcUtilsCommand)
}
