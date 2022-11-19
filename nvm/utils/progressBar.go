package utils

import (
	"fmt"
	"time"
)

// show a progress bar in stdout; needs bytes to be written using something like io.Reader.Read
func ProgressBar(b *[]byte, size int) {
	graph := "#"
	totalBlocks := 25
	lastBlocks := 0
	sleepTime := time.Millisecond * 100
	for {
		cur := len(*b)
		curBlocks := int(float32(cur) / float32(size) * float32(totalBlocks))

		if curBlocks > lastBlocks {
			// add another block to stdout progress
			bars := ""
			for len(bars) < curBlocks {
				bars += graph
			}

			// 25 needs to line up with `totalBlocks`
			fmt.Printf("\r[%-25s]", bars)
			lastBlocks = curBlocks
		}

		if cur == size {
			break
		}

		time.Sleep(sleepTime)
	}
	fmt.Println()
}
