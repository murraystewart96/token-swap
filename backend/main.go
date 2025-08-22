package main

import (
	"time"

	"github.com/murraystewart96/token-swap/cmd"
	"github.com/rs/zerolog"
)

var commitHash, version, buildTime string //nolint:gochecknoglobals

func main() {
	zerolog.TimeFieldFormat = time.RFC3339Nano

	cmd.Execute()
}
