package beater

import (
	"fmt"
	"os"
	"time"

	"github.com/KentaroAOKI/mssqlbeat/config"
)

// enabledArray filters the input array to include only enabled inputs.
func enabledArray(inputs []config.Input) []config.Input {
	var outs []config.Input
	for _, input := range inputs {
		if !input.Enabled {
			continue
		}
		outs = append(outs, input)
	}
	return outs
}

// chunkArray splits the input array into chunks of the specified size.
func chunkArray(inputs []config.Input, size int) [][]config.Input {
	var outs [][]config.Input
	for size < len(inputs) {
		inputs, outs = inputs[size:], append(outs, inputs[0:size:size])
	}
	outs = append(outs, inputs)
	return outs
}

// writeLastTime writes the given time to a file with the specified prefix.
func writeLastTime(file_prefix string, t time.Time) error {
	var lastTimeFile = fmt.Sprintf("last_time_%s.txt", file_prefix)
	return os.WriteFile(lastTimeFile, []byte(t.Format(time.RFC3339Nano)), 0644)
}

// readLastTime reads the last time from a file with the specified prefix.
// If the file does not exist, it creates the file with the specified time (current or oldest).
func readLastTime(file_prefix string, initializeWithCurrentTime bool) (time.Time, error) {
	var lastTimeFile = fmt.Sprintf("last_time_%s.txt", file_prefix)
	data, err := os.ReadFile(lastTimeFile)
	if err != nil {
		if os.IsNotExist(err) {
			var timeToWrite time.Time
			if initializeWithCurrentTime {
				timeToWrite = time.Now()
			} else {
				timeToWrite = time.Unix(0, 0) // Oldest possible time
			}
			writeLastTime(file_prefix, timeToWrite)
			return timeToWrite, nil
		}
		return time.Time{}, err
	}
	lastTime, err := time.Parse(time.RFC3339Nano, string(data))
	if err != nil {
		return time.Time{}, err
	}
	return lastTime, nil
}
