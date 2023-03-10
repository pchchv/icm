package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/op/go-logging"
)

const (
	size = 1024
)

var (
	Log    *Logger
	exited bool
	level  = logging.INFO // default level
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
)

type statusMsg struct {
	Text    string
	IsError bool
}

type Logger struct {
	*logging.Logger
	backend *logging.MemoryBackend
	logFile *os.File
	sLog    []statusMsg
}

func (c *Logger) FlushStatus() chan statusMsg {
	ch := make(chan statusMsg)

	go func() {
		for _, sm := range c.sLog {
			ch <- sm
		}

		close(ch)
		c.sLog = []statusMsg{}
	}()

	return ch
}

func (c *Logger) StatusQueued() bool {
	return len(c.sLog) > 0
}

func (c *Logger) Status(s string) {
	c.addStatus(statusMsg{s, false})
}

func (c *Logger) StatusErr(err error) {
	c.addStatus(statusMsg{err.Error(), true})
}

func (c *Logger) addStatus(sm statusMsg) {
	c.sLog = append(c.sLog, sm)
}

func (c *Logger) Statusf(s string, a ...interface{}) {
	c.Status(fmt.Sprintf(s, a...))
}

func Init() *Logger {
	if Log == nil {
		logging.SetFormatter(format) // setup default formatter

		Log = &Logger{
			logging.MustGetLogger("ctop"),
			logging.NewMemoryBackend(size),
			nil,
			[]statusMsg{},
		}

		debugMode := debugMode()
		if debugMode {
			level = logging.DEBUG
		}
		backendLvl := logging.AddModuleLevel(Log.backend)
		backendLvl.SetLevel(level, "")

		logFilePath := debugModeFile()
		if logFilePath == "" {
			logging.SetBackend(backendLvl)
		} else {
			logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o666)
			if err != nil {
				logging.SetBackend(backendLvl)
				Log.Error("Unable to create log file: %s", err.Error())
			} else {
				backendFile := logging.NewLogBackend(logFile, "", 0)
				backendFileLvl := logging.AddModuleLevel(backendFile)
				backendFileLvl.SetLevel(level, "")
				logging.SetBackend(backendLvl, backendFileLvl)
				Log.logFile = logFile
			}
		}

		if debugMode {
			StartServer()
		}

		Log.Notice("logger initialized")
	}

	return Log
}

func (log *Logger) tail() chan string {
	stream := make(chan string)

	node := log.backend.Head()
	go func() {
		for {
			stream <- node.Record.Formatted(0)
			for {
				nnode := node.Next()
				if nnode != nil {
					node = nnode
					break
				}
				if exited {
					close(stream)
					return
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return stream
}

func (log *Logger) Exit() {
	exited = true

	if log.logFile != nil {
		_ = log.logFile.Close()
	}

	StopServer()
}

func debugMode() bool {
	return os.Getenv("CTOP_DEBUG") == "1"
}

func debugModeTCP() bool {
	return os.Getenv("CTOP_DEBUG_TCP") == "1"
}

func debugModeFile() string {
	return os.Getenv("CTOP_DEBUG_FILE")
}
