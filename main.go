package gologs

import (
	"fmt"
	"os"
	"strings"
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
	CRITICAL_LVL
	SHOULDNT_HAPPEN_LVL
)

type GoLogger struct {
	PrintLogsLevel int

	LogLevel int
	LogFile  string
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

func (logger *GoLogger) leftPad(str string, padCount int, padStr rune) string {
	if len(str) >= padCount {
		return str
	}

	return strings.Repeat(string(padStr), padCount-len(str)) + str
}

// func (logger *GoLogger) rightPad(str string, padCount int, padStr rune) string {
// 	if len(str) >= padCount {
// 		return str
// 	}

// 	return str + strings.Repeat(string(padStr), padCount-len(str))
// }

func (logger *GoLogger) format_DDMMYYYY_HHMMSSMS(time time.Time) string {
	time = time.UTC()

	return fmt.Sprintf("%d/%s/%d %d:%d:%d.%d", time.Day(), logger.leftPad(fmt.Sprint(int(time.Month())), 2, '0'), time.Year(), time.Hour(), time.Minute(), time.Second(), time.UnixMilli()%SECOND)
}

func (logger *GoLogger) log(level int, str string, errs ...error) {
	var log_lvl_str = ""
	switch level {
	case DEBUG_LVL:
		log_lvl_str = "DEBUG"
	case INFO_LVL:
		log_lvl_str = "INFO"
	case IMPORTANT_LVL:
		log_lvl_str = "IMPORTANT"
	case WARN_LVL:
		log_lvl_str = "WARN"
	case ERROR_LVL:
		log_lvl_str = "ERROR"
	case CRITICAL_LVL:
		log_lvl_str = "CRITICAL"
	case SHOULDNT_HAPPEN_LVL:
		log_lvl_str = "SHOULDNT_HAPPEN"
	}

	str = fmt.Sprintf("[%s] %s: %s\n", log_lvl_str, logger.format_DDMMYYYY_HHMMSSMS(time.Now()), str)
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
	case CRITICAL_LVL:
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
