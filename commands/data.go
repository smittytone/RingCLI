package commands

import (
	"github.com/spf13/cobra"

	dataCommands "ringcli/commands/data"
	config "ringcli/lib/config"
)

// Define the `data` command.
var dataCommand = &cobra.Command{
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

	dataCommand.AddCommand(dataCommands.StepsCmd)
	dataCommand.AddCommand(dataCommands.HeartRateCmd)
	dataCommand.AddCommand(dataCommands.BloodOxygenCmd)
	dataCommand.AddCommand(dataCommands.SleepCmd)
	rootCmd.AddCommand(dataCommand)
}

// Display help when `data` or `utils` is called without args.
func showAppHelp(cmd *cobra.Command, args []string) {

	if config.Config.DoShowVersion {
		showVersion()
	} else {
		cmd.Help()
	}
}
