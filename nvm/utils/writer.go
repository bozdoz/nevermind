package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

var w = tabwriter.NewWriter(os.Stderr, 1, 1, 4, ' ', 0)

func PrintTabs(str string) {
	fmt.Fprintln(w, str)
}

func FlushTabs() {
	w.Flush()
}
