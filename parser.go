package main

import (
	"os"
	"time"
)

const FIRST_OFFSET = 0x6

// I'm using the offsets in https://en.wikipedia.org/wiki/MOI_(file_format)

func parseUint(sl []byte) uint {
	o := uint(0)

	for i, b := range sl {
		off := (len(sl) - i - 1) * 8

		o |= uint(b) << off
	}

	return o
}

func getTimestampAndDuration(f *os.File) (time.Time, time.Duration) {
	b := make([]byte, 0x11-FIRST_OFFSET+1)
	_, err := f.ReadAt(b[:], FIRST_OFFSET)
	if err != nil {
		panic("Can't read metadata: " + err.Error())
	}

	d := time.Date(
		int(parseUint(b[:2])),
		time.Month(b[2]),
		int(b[3]),
		int(b[4]),
		int(b[5]),
		0,
		int(b[6])*int(time.Millisecond/time.Nanosecond),
		time.UTC,
	)

	return d, time.Duration(parseUint(b[-FIRST_OFFSET+0xE : 0x11-FIRST_OFFSET+1])) * time.Millisecond
}
