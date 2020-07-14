package main

import (
	"fmt"
	"strings"
	"time"
)

func countup(chars int, seconds int) {
	var msecPerChar time.Duration
	msecPerChar = time.Duration(seconds * 1000 / chars)
	up := 0
	top := seconds * 1000
	nb := 0
	timer := time.Tick(msecPerChar * time.Millisecond)
	for up < top {
		<-timer
		// Dirty trick to overwrite the same line
		fmt.Printf("\r|%s%s|", strings.Repeat("=", nb), strings.Repeat(" ", chars-nb))

		up += int(msecPerChar)
		nb++
	}
	fmt.Printf("\n")
}
