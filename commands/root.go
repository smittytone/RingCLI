package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	rcLog "ringcli/lib/log"
)

type AppVersion struct {
	Major int
	Minor int
	Patch int
}

var (
	doShowVersion bool       = false
	Version       AppVersion = AppVersion{
		Major: 0,
		Minor: 0,
		Patch: 0,
	}
)

var rootCmd = &cobra.Command{
	Use:   "ringlci",
	Short: "A Colmi R02 CLI tool",
	Long:  "A Colmi R02 CLI tool.",
	Run:   showRootHelp,
}

func Execute() {

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func showRootHelp(cmd *cobra.Command, args []string) {

	if doShowVersion {
		rcLog.Report("RingCLI version %d.%d.%d\nCopyright © %d Tony Smith (@smittytone)", Version.Major, Version.Minor, Version.Patch, time.Now().Year())
	} else {
		showHelp()
	}
}

func init() {

	// Add persistent flags, ie. those spanning all commands and sub-commands.
	rootCmd.PersistentFlags().BoolVarP(&doShowVersion, "version", "", false, "Show tool version information")
}

func showHelp() {

	help := `ringcli
Provides ready access to a range of Microvisor Cloud and device functionality,

Usage:
  ringcli {COMMAND} [SUBCOMMAND] {REQUIRED VALUES} [FLAGS]

Commands:
  data       Access ring activity and heath data.
  utils      Scan for rings and get ring info, including battery state.

For more information on each command, run
  ringcli {COMMAND} --help
`
	fmt.Println(help)
}