package common

import (
	"io"
	"log"
	"os"
	"strings"
)

/*
Enables debugger with DEBUG=1

If you want node logs too, set DEBUG=*
*/
var DEBUG = os.Getenv("DEBUG")

/*
sets up logger if DEBUG=1

logging is disabled by default, but you could also
pass DEBUG=0 if that's something you think you need
to do
*/
func Debugger() {
	switch strings.ToLower(DEBUG) {
	case "", "false", "0", "n", "no":
		// falsy values discard logs
		log.SetOutput(io.Discard)
	}

	log.SetPrefix("\n[debug] ")

	// TODO: should flags be customized?
	log.SetFlags(0)
}
