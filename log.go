package main

import (
	"io"
	"os"
)

const EnvLog = "CLU_LOG"
const EnvLogFile = "CLU_LOG_PATH"

func logOutput() (logOutput io.Writer, err error) {
	logOutput = nil
	if os.Getenv(EnvLog) != "" {
		logOutput = os.Stderr

		if logPath := os.Getenv(EnvLogFile); logPath != "" {
			var err error
			logOutput, err = os.Create(logPath)
			if err != nil {
				return nil, err
			}
		}
	}
	return
}
