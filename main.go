package main

import (
	"fmt"
	"os"
	commands "ringcli/commands"
)

var (
	versionMajor int = 0
	versionMinor int = 1
	versionPatch int = 0
)

func main() {
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
