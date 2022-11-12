package utils

import (
	"fmt"
	"os"
	"text/tabwriter"
)

// customized writer to print help logs nicely
var Writer = tabwriter.NewWriter(os.Stderr, 1, 1, 4, ' ', 0)

// aligns tabs with tabwriter
func PrintTabs(str string) {
	fmt.Fprintln(Writer, str)
}

// flush tabbed output to stdout
func FlushTabs() {
	Writer.Flush()
}
