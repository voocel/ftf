package main

import (
	"fmt"
	"runtime"
)

var (
	Version = "v1.0.0"
	ID      = "N/A"
	Time    = "N/A"
)

const template =
`NAME
  ftf                 %s

GENERAL
  GOARCH:             %s
  GOOS:               %s
  Go Version:         %s
  Version:            %s
  Build ID:           %s
  Build Time:         %s

COMMANDS
  send                %s
  receive             %s
  help,h              %s
`

func Summary() string {
	return fmt.Sprintf(template,
		"a fast filetransfer for Golang",
		runtime.GOARCH,
		runtime.GOOS,
		runtime.Version(),
		Version,
		ID,
		Time,
		"send file(s), or folder",
		"receive file(s), or folder",
		"shows a list of commands or help for one command",
	)
}
