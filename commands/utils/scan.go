package rcUtilsCommands

import (
	"github.com/spf13/cobra"
	"os"
	ble "ringcli/lib/ble"
	errors "ringcli/lib/errors"
	log "ringcli/lib/log"
	"strings"
	"time"
	"tinygo.org/x/bluetooth"
	utils "ringcli/lib/utils"
)

const (
	SCAN_TIMEOUT_S = 30
)

type ScanRecord struct {
	Name string // BLE device local name
	Ad   []byte // BLE GAP ad packet
}

// Globals relevant only to this command
var (
	rings            map[string]string     = make(map[string]string)     // Dictionary of rings. Key is BLE address
	devices          map[string]ScanRecord = make(map[string]ScanRecord) // Dictionary of other devices. Key is BLE address. Debug only
	scanTimer        *time.Timer                                         // Scan window timer
	scanForFirstRing bool                  = false                       // Stop scanning on first ring found
	doBind           bool                  = false                       // Auto-bind ring
)

// Define the `scan` subcommand.
var ScanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan for rings",
	Long:  "Scan for rings",
	Run:   doScan,
}

func doScan(cmd *cobra.Command, args []string) {

	log.Prompt("Scanning for rings")

	// Enable BLE
	radio := ble.Open()

	// Establish a timer for the scan
	scanTimer = time.NewTimer(SCAN_TIMEOUT_S * time.Second)
	defer scanTimer.Stop()
	go func() {
		<-scanTimer.C
		// Timeout
		exitCode := errors.ERROR_CODE_SCAN_TIMEOUT

		log.ClearPrompt()
		if len(rings) > 0 {
			// We have one or more rings, so display their data before exiting
			printFoundRings()
			exitCode = errors.ERROR_CODE_NONE
		} else {
			log.ReportError("Scan timed out and no rings found")
		}

		// Display de-duped list of other BLE devices on debug runs
		if debug && len(devices) > 0 {
			for address, record := range devices {
				log.Report("Device %s with BLE address %s (%v)", record.Name, address, record.Ad)
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
			log.ClearPrompt()
			printFoundRings()
			os.Exit(errors.ERROR_CODE_NONE)
		}
	} else if debug {
		// Only note other devices on debug runs
		if _, ok := devices[address]; !ok {
			record := ScanRecord{}
			if name != "" {
				record.Name = name
			} else {
				record.Name = "Unnamed"
			}

			record.Ad = make([]byte, 0, 31)
			if device.AdvertisementPayload != nil {
				record.Ad = append(record.Ad, device.AdvertisementPayload.Bytes()...)
			} else {
				record.Ad = append(record.Ad, []byte{0x01, 0x02, 0x03, 0x04}...)
			}

			devices[address] = record
		}
	}
}

func printFoundRings() {

	bound := false
	header := "Ring found:"
	if len(rings) > 1 {
		log.Report("Rings found:")
		header = " "
	}

	for address, name := range rings {
		log.Report("%s %s with BLE address %s", header, name, address)
		if doBind && !bound {
			utils.MakeBinding(address, name, true)
			log.Report("Ring %s bound", address)
			bound = true
		}
	}
}
