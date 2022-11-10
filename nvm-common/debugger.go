package common

import (
	"io"
	"log"
	"os"
	"strings"
)

func Debugger() {
	DEBUG := os.Getenv("DEBUG")

	switch strings.ToLower(DEBUG) {
	case "", "false", "0", "n", "no":
		// falsy values discard logs
		log.SetOutput(io.Discard)
	}

	log.SetPrefix("\n[debug] ")

	// TODO: should flags be customized?
	log.SetFlags(0)
}
