package logger

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//ESLogger defined
type ESLogger struct {
	Log *logrus.Logger
}

type LogInfo struct {
	method     string
	path       string
	clientIP   string
	userAgent  string
	dataLength int
	function   string
	line       int
}

//Logger interface
type Logger interface {
	Debug(direction string, i *LogInfo, msg string)
	Info(direction string, i *LogInfo, msg string)
	Error(direction string, i *LogInfo, msg string)
	Fatal(direction string, i *LogInfo, msg string)
	GetLogger() *logrus.Logger
}

//InitLogger used to get Logger componnet
func InitLogger(logLevel int) *ESLogger {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	log := logrus.New()
	log.Level = logrus.DebugLevel

	return &ESLogger{
		log}
}

func (logger *ESLogger) GetLogger() *logrus.Logger {
	return logger.Log
}

func (logger *ESLogger) Debug(direction string, i *LogInfo, msg string) {
	logger.Log.WithFields(logrus.Fields{
		"method":     i.method,
		"path":       i.path,
		"direction":  direction,
		"clientIP":   i.clientIP,
		"userAgent":  i.userAgent,
		"dataLength": i.dataLength,
		"func":       i.function,
		"line":       i.line,
	}).Debug(msg)
}

func (logger *ESLogger) Info(direction string, i *LogInfo, msg string) {
	logger.Log.WithFields(logrus.Fields{
		"method":     i.method,
		"path":       i.path,
		"direction":  direction,
		"clientIP":   i.clientIP,
		"userAgent":  i.userAgent,
		"dataLength": i.dataLength,
		"func":       i.function,
		"line":       i.line,
	}).Info(msg)
}

func (logger *ESLogger) Error(direction string, i *LogInfo, msg string) {
	logger.Log.WithFields(logrus.Fields{
		"method":     i.method,
		"path":       i.path,
		"direction":  direction,
		"clientIP":   i.clientIP,
		"userAgent":  i.userAgent,
		"dataLength": i.dataLength,
		"func":       i.function,
		"line":       i.line,
	}).Error(msg)
}

func (logger *ESLogger) Fatal(direction string, i *LogInfo, msg string) {
	logger.Log.WithFields(logrus.Fields{
		"method":     i.method,
		"path":       i.path,
		"direction":  direction,
		"clientIP":   i.clientIP,
		"userAgent":  i.userAgent,
		"dataLength": i.dataLength,
		"func":       i.function,
		"line":       i.line,
	}).Fatal(msg)
}

func Log(logger Logger, level string, direction string, i *LogInfo, msg string) {
	switch level {
	case "Debug":
		logger.Debug(direction, i, msg)
	case "Info":
		logger.Info(direction, i, msg)
	case "Error":
		logger.Error(direction, i, msg)
	case "Fatal":
		logger.Fatal(direction, i, msg)
	default:
		fmt.Printf("no such log level %s.", level)
	}
}

func Trace() *LogInfo {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	_, line := f.FileLine(pc[0])
	return &LogInfo{
		function: f.Name(),
		line:     line,
	}
}

func BuildLogInfo(c *gin.Context) *LogInfo {

	dataLength := c.Writer.Size()
	if dataLength < 0 {
		dataLength = 0
	}
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	_, line := f.FileLine(pc[0])

	return &LogInfo{
		method:     c.Request.Method,
		path:       c.Request.URL.Path,
		clientIP:   c.ClientIP(),
		userAgent:  c.Request.UserAgent(),
		dataLength: dataLength,
		function:   f.Name(),
		line:       line,
	}
}
