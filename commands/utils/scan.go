package rcUtilsCommands

import (
	"os"
	"time"
	"strings"
	// External code
	"github.com/spf13/cobra"
	"tinygo.org/x/bluetooth"
	// Library code
	rcBLE "ringcli/lib/ble"
	rcErrors "ringcli/lib/errors"
	rcLog "ringcli/lib/log"
)

const (
	SCAN_TIMEOUT_S = 30
)

// Globals relevant only to this command
var (
	rings            map[string]string   = make(map[string]string)
	devices          map[string]string = make(map[string]string)
	scanTimer        *time.Timer
	bspCount         int
	scanForFirstRing bool = false
)

// Define the `scan` subcommand.
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for rings",
	Long:  "Scan for rings",
	Run:   doScan,
}

func doScan(cmd *cobra.Command, args []string) {

	bspCount = rcLog.Raw("Scanning...")

	// Enable BLE
	ble := rcBLE.Open()

	// Establish a timer for the scan
	scanTimer = time.NewTimer(SCAN_TIMEOUT_S * time.Second)
	defer scanTimer.Stop()
	go func() {
		<-scanTimer.C
		exitCode := rcErrors.ERROR_CODE_SCAN_TIMEOUT
		if len(rings) > 0 {
			// We have one or more rings, so display their data
			printFoundRings()
			exitCode = rcErrors.ERROR_CODE_NONE
		} else {
			rcLog.ReportError("Scan timed out and no rings found")
		}

		// Display de-duped list of other BLE devices on debug runs
		if debug && len(devices) > 0 {
			for address, name := range devices {
				rcLog.Report("Device %s with BLE address $s", name, address)
			}
		}

		os.Exit(exitCode)
	}()

	// Setup done, so initiate a scan
	rcBLE.BeginScan(ble, onScan)
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
			printFoundRings()
			os.Exit(rcErrors.ERROR_CODE_NONE)
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

	rcLog.Backspaces(bspCount)
	if len(rings) > 0 {
		for address, name := range rings {
			rcLog.Report("Ring found: %s with BLE address %s", name, address)
		}
	}
}