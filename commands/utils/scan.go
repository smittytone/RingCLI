package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	"os"
	ble "ringcli/lib/ble"
	errors "ringcli/lib/errors"
	log "ringcli/lib/log"
	utils "ringcli/lib/utils"
	"strings"
	"time"
	"tinygo.org/x/bluetooth"
)

const (
	SCAN_TIMEOUT_S = 30
)

// Globals relevant only to this command
var (
	rings            map[string]string = make(map[string]string) // Dictionary of rings. Key is BLE address
	devices          map[string]string = make(map[string]string) // Dictionary of other devices. Key is BLE address. Debug only
	scanTimer        *time.Timer                                 // Scan window timer
	scanForFirstRing bool              = false                   // Stop scanning on first ring found
)

// Define the `scan` subcommand.
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for rings",
	Long:  "Scan for rings",
	Run:   doScan,
}

func doScan(cmd *cobra.Command, args []string) {

	bspCount = log.Raw("Scanning...  ")
	utils.AnimateCursor()

	// Enable BLE
	radio := ble.Open()

	// Establish a timer for the scan
	scanTimer = time.NewTimer(SCAN_TIMEOUT_S * time.Second)
	defer scanTimer.Stop()
	go func() {
		<-scanTimer.C
		utils.StopAnimation()
		exitCode := errors.ERROR_CODE_SCAN_TIMEOUT
		if len(rings) > 0 {
			// We have one or more rings, so display their data before exiting
			printFoundRings()
			exitCode = errors.ERROR_CODE_NONE
		} else {
			log.ReportError("Scan timed out and no rings found")
		}

		// Display de-duped list of other BLE devices on debug runs
		if debug && len(devices) > 0 {
			for address, name := range devices {
				log.Report("Device %s with BLE address $s", name, address)
			}
		}

		os.Exit(exitCode)
	}()

	// Setup done, so initiate a scan
	ble.BeginScan(radio, onScan)
}

func onScan(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {

	address := device.Address.String()
	name := device.LocalName()

	// Only proceed when we've found a ring
	// TODO Check other known prefixes too
	if strings.HasPrefix(name, "R02_") {
		if rings[address] == "" {
			rings[address] = name
		}

		if scanForFirstRing {
			scanTimer.Stop()
			utils.StopAnimation()
			printFoundRings()
			os.Exit(errors.ERROR_CODE_NONE)
		}
	} else if debug {
		// Only note other devices on debug runs
		if devices[address] == "" {
			if name != "" {
				devices[address] = name
			} else {
				devices[address] = "Unnamed"
			}
		}
	}
}

func printFoundRings() {

	log.Backspaces(bspCount)

	if len(rings) > 0 {
		for address, name := range rings {
			log.Report("Ring found: %s with BLE address %s", name, address)
		}
	}
}
