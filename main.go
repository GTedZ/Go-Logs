package gologs

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/sys/windows"
)

const (
	SECOND = 1000

	ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004

	COLOR_RESET  = "\x1b[0m"
	COLOR_RED    = "\x1b[31m"
	COLOR_GREEN  = "\x1b[32m"
	COLOR_YELLOW = "\x1b[33m"
	COLOR_BLUE   = "\x1b[34m"
	COLOR_PURPLE = "\x1b[35m"
	COLOR_CYAN   = "\x1b[36m"
	COLOR_WHITE  = "\x1b[37m"
	COLOR_GRAY   = "\x1b[90m"
)

func init() {
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

const (
	DEBUG_LVL = iota
	INFO_LVL
	IMPORTANT_LVL
	WARN_LVL
	ERROR_LVL
	SHOULDNT_HAPPEN_LVL

	NOTHING_LVL
)

type GoLogger struct {
	disabled bool

	PrintLogsLevel int

	LogLevel int
	LogFile  string

	OnDebug          func(logStr string, message string, errs ...error)
	OnInfo           func(logStr string, message string, errs ...error)
	OnImportant      func(logStr string, message string, errs ...error)
	OnWarn           func(logStr string, message string, errs ...error)
	OnError          func(logStr string, message string, errs ...error)
	OnShouldntHappen func(logStr string, message string, errs ...error)
}

func (logger *GoLogger) DEBUG(message string, err ...error) {
	logger.log(DEBUG_LVL, message, err...)
}

func (logger *GoLogger) INFO(message string, err ...error) {
	logger.log(INFO_LVL, message, err...)
}

func (logger *GoLogger) IMPORTANT(message string, err ...error) {
	logger.log(IMPORTANT_LVL, message, err...)
}

func (logger *GoLogger) WARN(message string, err ...error) {
	logger.log(WARN_LVL, message, err...)
}

func (logger *GoLogger) ERROR(message string, err ...error) {
	logger.log(ERROR_LVL, message, err...)
}

func (logger *GoLogger) SHOULDNT_HAPPEN(message string, err ...error) {
	logger.log(SHOULDNT_HAPPEN_LVL, message, err...)
}

func (logger *GoLogger) format_DDMMYYYY_HHMMSSMS(time time.Time) string {
	time = time.UTC()

	return fmt.Sprintf("%02d/%02d/%d %02d:%02d:%02d.%03d", time.Day(), int(time.Month()), time.Year(), time.Hour(), time.Minute(), time.Second(), time.UnixMilli()%SECOND)
}

func (logger *GoLogger) Enable() {
	logger.disabled = true
}

func (logger *GoLogger) Disable() {
	logger.disabled = false
}

func (logger *GoLogger) log(level int, message string, errs ...error) {
	if logger.disabled {
		return
	}

	var cb_func func(logStr string, message string, errs ...error)

	var log_lvl_str = ""
	switch level {
	case DEBUG_LVL:
		log_lvl_str = "DEBUG"
		cb_func = logger.OnDebug
	case INFO_LVL:
		log_lvl_str = "INFO"
		cb_func = logger.OnInfo
	case IMPORTANT_LVL:
		log_lvl_str = "IMPORTANT"
		cb_func = logger.OnImportant
	case WARN_LVL:
		log_lvl_str = "WARN"
		cb_func = logger.OnWarn
	case ERROR_LVL:
		log_lvl_str = "ERROR"
		cb_func = logger.OnError
	case SHOULDNT_HAPPEN_LVL:
		log_lvl_str = "SHOULDNT_HAPPEN"
		cb_func = logger.OnShouldntHappen
	}

	str := fmt.Sprintf("[%s] %s: %s\n", log_lvl_str, logger.format_DDMMYYYY_HHMMSSMS(time.Now()), message)
	for i, err := range errs {
		errIndex_str := ""
		if len(errs) > 1 {
			errIndex_str = fmt.Sprintf(" %d", i+1)
		}
		str += fmt.Sprintf("\t\t\t\t=> err%s: %s\n", errIndex_str, err.Error())
	}

	var Color string
	switch level {
	case DEBUG_LVL:
		Color = COLOR_GRAY
	case INFO_LVL:
		Color = COLOR_WHITE
	case IMPORTANT_LVL:
		Color = COLOR_GREEN
	case WARN_LVL:
		Color = COLOR_YELLOW
	case ERROR_LVL:
		Color = COLOR_RED
	case SHOULDNT_HAPPEN_LVL:
		Color = COLOR_RED
	}

	if logger.PrintLogsLevel <= level {
		fmt.Println(Color + str + COLOR_RESET)
	}

	if logger.LogLevel <= level {
		logger.appendLogFile(str)
	}

	if cb_func != nil {
		cb_func(str, message, errs...)
	}
}

func (logger *GoLogger) appendLogFile(str string) {
	fileName := "logs.txt"
	if logger.LogFile != "" {
		fileName = logger.LogFile
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("[GoLogger] Error opening file:", err)
		return
	}
	defer file.Close()

	// Write data to file
	_, err = file.WriteString(str)
	if err != nil {
		fmt.Println("[GoLogger] Error writing to file:", err)
		return
	}
}
