package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"

	Log "ringcli/lib/log"
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
	ValidArgs: []string{
		"data",
		"utils",
	},
}

func Execute() {

	Log.CursorHide()
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
	Log.CursorShow()
}

func showRootHelp(cmd *cobra.Command, args []string) {

	if doShowVersion {
		Log.Report("RingCLI version %d.%d.%d\nCopyright © %d Tony Smith (@smittytone)", Version.Major, Version.Minor, Version.Patch, time.Now().Year())
	} else {
		showHelp()
	}
}

func init() {

	// Add persistent flags, ie. those spanning all commands and sub-commands.
	//rootCmd.PersistentFlags().BoolVarP(&doShowVersion, "version", "", false, "Show tool version information")
}

func showHelp() {

	help := `Manage your Colmi R02 smart ring and retrieve data from it.

Usage:
  ringcli {COMMAND} [SUB-COMMAND] {REQUIRED VALUES} [FLAGS]

Commands:
  data       Access ring activity (steps, calories burned, distance moved) and heath data,
             including daily heart rate, blood oxygen (SpO2) and sleep records.
  utils      Scan for nearby rings, get specific ring info, including battery state,
             and perform housekeeping tasks such as flashing a ring's LED to help find it,
             bind the ring's address to a local store to save typing, set and enable periodic
             heart rate readings, set the ring's internal clock to (re)initialise the ring, and
             shut the ring down.

For more information on each command's array of sub-commands, run 'ringcli {COMMAND} --help'

All sub-commands other than 'ringcli utils scan' require the target ring's BLE address.
The 'scan' sub-command will get this for you. Use the 'bind' sub-command to retain this value
locally so you need not enter it again unless you have multiple rings. Only one ring can be
bound to a local machine. This is not a pairing process: it is simply a convenience feature.
`
	fmt.Println(help)
}
