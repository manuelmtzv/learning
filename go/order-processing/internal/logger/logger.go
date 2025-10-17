package logger

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

func New() *log.Logger {
	logger := log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    true,
		ReportTimestamp: true,
		TimeFormat:      time.StampMicro,
	})

	return logger
}
