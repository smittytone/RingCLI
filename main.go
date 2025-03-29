package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	commands "ringcli/commands"
	rcBLE "ringcli/lib/ble"
)

var (
	versionMajor int = 0
	versionMinor int = 1
	versionPatch int = 0
)

func main() {
	signalChannel := make(chan os.Signal, 2)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		rcBLE.Clean()
		os.Exit(0)
	}()

	commands.Version.Major = versionMajor
	commands.Version.Minor = versionMinor
	commands.Version.Patch = versionPatch

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r)
			os.Exit(1)
		}
	}()
	commands.Execute()
}
