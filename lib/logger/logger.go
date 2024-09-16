package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Settings struct {
	Path       string `yaml:"path"`
	Name       string `yaml:"name"`
	Ext        string `yaml:"ext"`
	TimeFormat string `yaml:"time_format"`
}

type loglevel int

const (
	DEBUG loglevel = iota
	INFO
	WARING
	ERROR
	FATAL
)

const (
	flags              = log.LstdFlags
	defaultCallerDepth = 2
	bufferSize         = 1e5
)

type logEntry struct {
	msg   string
	level loglevel
}

var (
	levelFlags = []string{"DEBUG", "INFO", "WARING", "ERROR", "FATAL"}
)

type Logger struct {
	logFile   *os.File
	logger    *log.Logger
	entryChan chan *logEntry
	entryPool *sync.Pool
}

var DefaultLogger = MakeStdoutLogger()

// MakeStdoutLogger creates a logger which print msg to stdout
func MakeStdoutLogger() *Logger {
	logger := &Logger{
		logFile:   nil,
		logger:    log.New(os.Stdout, "", flags),
		entryChan: make(chan *logEntry, bufferSize),
		entryPool: &sync.Pool{
			New: func() interface{} {
				return &logEntry{}
			},
		},
	}
	go func() {
		for e := range logger.entryChan {
			_ = logger.logger.Output(0, e.msg)
			logger.entryPool.Put(e)
		}
	}()
	return logger
}

// NewFileLogger crates a Logger which print msg to stdout and log file
func NewFileLogger(settings *Settings) (*Logger, error) {
	fileName := fmt.Sprintf("%s-%s.%s",
		settings.Name,
		time.Now().Format(settings.TimeFormat),
		settings.Ext)
	logFile, err := mustOpen(fileName, settings.Path)
	if err != nil {
		return nil, fmt.Errorf("cannot open log file: %v", err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	logger := &Logger{
		logFile:   logFile,
		logger:    log.New(mw, "", flags),
		entryChan: make(chan *logEntry, bufferSize),
		entryPool: &sync.Pool{
			New: func() interface{} {
				return &logEntry{}
			},
		},
	}
	go func() {
		for e := range logger.entryChan {
			logFileName := fmt.Sprintf("%s-%s.%s",
				settings.Name,
				time.Now().Format(settings.TimeFormat),
				settings.Ext)
			if path.Join(settings.Path, logFileName) != logger.logFile.Name() {
				logFile, err := mustOpen(logFileName, settings.Path)
				if err != nil {
					panic("open log " + logFileName + " failed: " + err.Error())
				}
				logger.logFile = logFile
				logger.logger = log.New(io.MultiWriter(os.Stdout, logFile), "", flags)
			}
			_ = logger.logger.Output(0, e.msg)
			logger.entryPool.Put(e)
		}
	}()
	return logger, nil
}

// Setup init DefaultLogger
func Setup(settings *Settings) {
	logger, err := NewFileLogger(settings)
	if err != nil {
		panic(err)
	}
	DefaultLogger = logger
}

// Output sends a msg to logger
func (l *Logger) Output(level loglevel, callerDepth int, msg string) {
	var formattedMsg string
	_, file, line, ok := runtime.Caller(callerDepth)
	if ok {
		formattedMsg = fmt.Sprintf("[%s][%s:%d] %s", levelFlags[level], filepath.Base(file), line, msg)
	} else {
		formattedMsg = fmt.Sprintf("[%s] %s", levelFlags[level], msg)
	}
	entry := l.entryPool.Get().(*logEntry)
	entry.level = level
	entry.msg = formattedMsg
	l.entryChan <- entry
}

func Debug(v ...interface{}) {
	msg := fmt.Sprint(v...)
	DefaultLogger.Output(DEBUG, defaultCallerDepth, msg)
}

// Debugf logs debug msg through DefaultLogger with format
func Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	DefaultLogger.Output(DEBUG, defaultCallerDepth, msg)
}

func Info(v ...interface{}) {
	msg := fmt.Sprint(v...)
	DefaultLogger.Output(INFO, defaultCallerDepth, msg)
}

func Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	DefaultLogger.Output(INFO, defaultCallerDepth, msg)
}

func Warn(v ...interface{}) {
	msg := fmt.Sprint(v...)
	DefaultLogger.Output(WARING, defaultCallerDepth, msg)
}

func Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	DefaultLogger.Output(WARING, defaultCallerDepth, msg)
}

func Error(v ...interface{}) {
	msg := fmt.Sprint(v...)
	DefaultLogger.Output(ERROR, defaultCallerDepth, msg)
}

func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	DefaultLogger.Output(ERROR, defaultCallerDepth, msg)
}

// Fatal prints error message then stop the program
func Fatal(v ...interface{}) {
	msg := fmt.Sprint(v...)
	DefaultLogger.Output(FATAL, defaultCallerDepth, msg)
}

//TODO:here Fatal does not stop the program, I do not know when it will be stoped
