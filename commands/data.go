package commands

import (
	"github.com/spf13/cobra"

	rcDataCommands "ringcli/commands/data"
)

// Define the `data` command.
var rcDataCommand = &cobra.Command{
	Use:   "data",
	Short: "Access ring data",
	Long:  "Read and store data retrieved from the ring.",
	Run:   showAppHelp,
	ValidArgs: []string{
		"heart",
		"steps",
	},
}

func init() {

	rcDataCommand.AddCommand(rcDataCommands.StepsCmd)
	rcDataCommand.AddCommand(rcDataCommands.HeartRateCmd)
	rcDataCommand.AddCommand(rcDataCommands.BloodOxygenCmd)
	rcDataCommand.AddCommand(rcDataCommands.SleepCmd)
	rootCmd.AddCommand(rcDataCommand)
}

// Display help when `app` is called without args.
func showAppHelp(cmd *cobra.Command, args []string) {

	cmd.Help()
}
