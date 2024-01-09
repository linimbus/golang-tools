package main

import (
	"io"
	"os"
	"time"
)

func main() {
	var latestBuffer []byte
	for {
		tempBuffer := make([]byte, 4096)
		cnt, err := os.Stdin.Read(tempBuffer[:])
		if err != nil && err == io.EOF {
			return
		}

		latestBuffer = append(latestBuffer, tempBuffer[:cnt]...)
		beginIndex := 0
		for i, char := range latestBuffer[:] {
			if char == '\n' || char == '\r' {
				os.Stdout.WriteString(time.Now().Format("2006-01-02 15:04:05.000") + " ")
				os.Stdout.Write(latestBuffer[beginIndex : i+1])
				beginIndex = i + 1
			}
		}
		if beginIndex != 0 {
			latestBuffer = latestBuffer[beginIndex:]
		}
	}
}
