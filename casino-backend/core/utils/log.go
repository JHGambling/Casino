package utils

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// LogLevel represents a log level configuration
type LogLevel struct {
	Char      string
	BoldColor *color.Color
	Plain     *color.Color
}

// Define log levels
var logLevels = map[string]LogLevel{
	"info":  {"INF", color.New(color.FgWhite), color.New(color.FgWhite)},
	"ok":    {"OK ", color.New(color.FgHiGreen), color.New(color.FgWhite)},
	"warn":  {"WRN", color.New(color.FgHiYellow), color.New(color.FgWhite)},
	"error": {"ERR", color.New(color.FgHiRed), color.New(color.FgWhite)},
	"fatal": {"FTL", color.New(color.FgHiRed), color.New(color.FgWhite)},
	"debug": {"DBG", color.New(color.FgMagenta), color.New(color.FgWhite)},
}

// Log prints a formatted log message
func Log(level string, origin string, v ...any) {
	lvl, ok := logLevels[level]
	if !ok {
		lvl = logLevels["info"]
	}

	// Time in light gray
	timestamp := time.Now().Format("15:04:05.000")
	timeColor := color.New(color.FgHiBlack)

	// Output
	timeColor.Printf("%s ", timestamp)
	lvl.BoldColor.Printf("%s ", lvl.Char)
	lvl.Plain.Printf("%s: ", origin)
	fmt.Println(fmt.Sprint(v...))
}
