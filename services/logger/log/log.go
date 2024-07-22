package log

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/labstack/echo"
	"sync"

	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// Log global
var Log *logrus.Logger
var singletonLogger = &MyLogger{}

/*func init() {
	InitLogger(false)
}*/

// MyLogger extend logrus.MyLogger
type MyLogger struct {
	*logrus.Logger
}

// Logger returns singletonLogger
// NOTE: Must run InitLogger() before
func Logger() *MyLogger {
	return singletonLogger
}

// InitLogger return singleton logger
func InitLogger(forTest bool) *MyLogger {
	if Log != nil {
		singletonLogger = &MyLogger{
			Logger: Log,
		}
		return singletonLogger
	}

	Log = logrus.New()

	if !forTest {
		writerInfo, err := rotatelogs.New(
			"./log_files/info/"+os.Getenv("APP_NAME")+"_%Y%m%d_info.log",
			rotatelogs.WithMaxAge(30*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			log.Printf("Failed to create rotatelogs: %s", err)
			return nil
		}

		writerError, err := rotatelogs.New(
			"./log_files/errors/"+os.Getenv("APP_NAME")+"_%Y%m%d_error.log",
			rotatelogs.WithMaxAge(30*24*time.Hour),
			rotatelogs.WithRotationTime(24*time.Hour),
		)
		if err != nil {
			log.Printf("Failed to create rotatelogs: %s", err)
			return nil
		}

		Log.Hooks.Add(lfshook.NewHook(
			lfshook.WriterMap{
				logrus.ErrorLevel: writerError,
				logrus.WarnLevel:  writerInfo,
				logrus.InfoLevel:  writerInfo,
			},
			&logrus.TextFormatter{
				TimestampFormat:  time.RFC3339Nano,
				QuoteEmptyFields: true,
			},
		))
	}

	singletonLogger = &MyLogger{
		Logger: Log,
	}
	return singletonLogger
}

func LoggerHandler(next gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := c.Request // change from ctx.Request to c.Request
		if !strings.Contains(req.RequestURI, "healthcheck") {
			// add some default fields to the logger ~ on all messages
			logger := Log.WithFields(logrus.Fields{
				"id":          req.Header.Get(echo.HeaderXRequestID),
				"method":      req.Method,
				"ip":          c.ClientIP(), // change from ctx.RealIP() to c.ClientIP()
				"request_uri": req.RequestURI,
			})
			c.Set("logger", logger) // change from ctx.Set to c.Set
			startTime := time.Now()

			defer func() {
				rsp := c.Writer.Status() // change from ctx.Response() to c.Writer.Status()
				// at the end we will want to log a few more interesting fields
				logger.WithFields(logrus.Fields{
					"status_code":  rsp,
					"runtime_nano": time.Since(startTime).Nanoseconds(),
				}).Info("Finished request")
			}()

			// now we will log out that we have actually started the request
			logger.WithFields(logrus.Fields{
				"id":             req.Header.Get(echo.HeaderXRequestID),
				"user_agent":     req.UserAgent(),
				"content_length": req.ContentLength,
			}).Info("Starting request")
		}

		next(c) // change from next(ctx) to next(c)

	}
}

// -----------------------------------
// Logger uses for trace
// -----------------------------------

// Args output message of print level
func Args(mess string, args ...interface{}) {
	a, _ := json.Marshal(args)
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"args": string(a),
		"file": file,
		"line": line,
	}).Info(mess)
}

// Print output message of print level
func Print(i ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Print(i...)
}

// Printf output format message of print level
func Printf(format string, i ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Printf(format, i...)
}

// Debug output message of debug level
func Debug(i ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Debug(i...)
}

// Debugf output format message of debug level
func Debugf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Debugf(format, args...)
}

// Info output message of info level
func Info(i ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Info(i...)
}

// Infof output format message of info level
func Infof(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Infof(format, args...)
}

// Warn output message of warn level
func Warn(i ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Warn(i...)
}

// Warnf output format message of warn level
func Warnf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Warnf(format, args...)
}

// Error output message of errors level
func Error(i ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Error(i...)
}

// Errorf output format message of errors level
func Errorf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Errorf(format, args...)
}

// Fatal output message of fatal level
func Fatal(i ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Fatal(i...)
}

// Fatalf output format message of fatal level
func Fatalf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Fatalf(format, args...)
}

// Panic output message of panic level
func Panic(i ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Panic(i...)
}

// Panicf output format message of panic level
func Panicf(format string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	singletonLogger.WithFields(logrus.Fields{
		"file": file,
		"line": line,
	}).Panicf(format, args...)
}

// To logrus.Level
func toLogrusLevel(level log.Lvl) logrus.Level {
	switch level {
	case log.DEBUG:
		return logrus.DebugLevel
	case log.INFO:
		return logrus.InfoLevel
	case log.WARN:
		return logrus.WarnLevel
	case log.ERROR:
		return logrus.ErrorLevel
	}

	return logrus.InfoLevel
}

// To Echo.log.lvl
func toEchoLevel(level logrus.Level) log.Lvl {
	switch level {
	case logrus.DebugLevel:
		return log.DEBUG
	case logrus.InfoLevel:
		return log.INFO
	case logrus.WarnLevel:
		return log.WARN
	case logrus.ErrorLevel:
		return log.ERROR
	}

	return log.OFF
}

// Output return logger io.Writer
func (l *MyLogger) Output() io.Writer {
	return l.Out
}

// SetOutput logger io.Writer
func (l *MyLogger) SetOutput(w io.Writer) {
	l.Out = w
}

// Level return logger level
func (l *MyLogger) Level() log.Lvl {
	return toEchoLevel(l.Logger.Level)
}

// SetLevel logger level
func (l *MyLogger) SetLevel(v log.Lvl) {
	l.Logger.Level = toLogrusLevel(v)
}

// Formatter return logger formatter
func (l *MyLogger) Formatter() logrus.Formatter {
	return l.Logger.Formatter
}

// SetFormatter logger formatter
// Only support logrus formatter
func (l *MyLogger) SetFormatter(formatter logrus.Formatter) {
	l.Logger.Formatter = formatter
}

// Prefix return logger prefix
// This function do nothing
func (l *MyLogger) Prefix() string {
	return ""
}

// SetPrefix logger prefix
// This function do nothing
func (l *MyLogger) SetPrefix(p string) {
	// do nothing
}

// -----------------------------------
// Logger uses for Echo
// -----------------------------------

// Print output message of print level
func (l *MyLogger) Print(i ...interface{}) {
	l.Logger.Print(i...)
}

// Printf output format message of print level
func (l *MyLogger) Printf(format string, args ...interface{}) {
	l.Logger.Printf(format, args...)
}

// Printj output json of print level
func (l *MyLogger) Printj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Println(string(b))
}

// Debug output message of debug level
func (l *MyLogger) Debug(i ...interface{}) {
	l.Logger.Info(i...)
}

// Debugf output format message of debug level
func (l *MyLogger) Debugf(format string, args ...interface{}) {
	l.Logger.Debugf(format, args...)
}

// Debugj output message of debug level
func (l *MyLogger) Debugj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Debugln(string(b))
}

// Info output message of info level
func (l *MyLogger) Info(i ...interface{}) {
	l.Logger.Info(i...)
}

// Infof output format message of info level
func (l *MyLogger) Infof(format string, args ...interface{}) {
	l.Logger.Infof(format, args...)
}

// Infoj output json of info level
func (l *MyLogger) Infoj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Infoln(string(b))
}

// Warn output message of warn level
func (l *MyLogger) Warn(i ...interface{}) {
	l.Logger.Warn(i...)
}

// Warnf output format message of warn level
func (l *MyLogger) Warnf(format string, args ...interface{}) {
	l.Logger.Warnf(format, args...)
}

// Warnj output json of warn level
func (l *MyLogger) Warnj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Warnln(string(b))
}

// Error output message of errors level
func (l *MyLogger) Error(i ...interface{}) {
	l.Logger.Error(i...)
}

// Errorf output format message of errors level
func (l *MyLogger) Errorf(format string, args ...interface{}) {
	l.Logger.Errorf(format, args...)
}

// Errorj output json of errors level
func (l *MyLogger) Errorj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Errorln(string(b))
}

// Fatal output message of fatal level
func (l *MyLogger) Fatal(i ...interface{}) {
	l.Logger.Fatal(i...)
}

// Fatalf output format message of fatal level
func (l *MyLogger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatalf(format, args...)
}

// Fatalj output json of fatal level
func (l *MyLogger) Fatalj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Fatalln(string(b))
}

// Panic output message of panic level
func (l *MyLogger) Panic(i ...interface{}) {
	l.Logger.Panic(i...)
}

// Panicf output format message of panic level
func (l *MyLogger) Panicf(format string, args ...interface{}) {
	l.Logger.Panicf(format, args...)
}

// Panicj output json of panic level
func (l *MyLogger) Panicj(j log.JSON) {
	b, err := json.Marshal(j)
	if err != nil {
		panic(err)
	}
	l.Logger.Panicln(string(b))
}

func (l *MyLogger) SetHeader(h string) {
	l.Logger.Info("Not implement yet")
}

// -----------------------------------
var (
	DefaultLogger    *logrus.Logger
	DefaultBaseEntry *logrus.Entry
	initOnce         sync.Once
)

func Init(name string) {
	initOnce.Do(func() {
		DefaultLogger = logrus.New()
		if l, e := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); e == nil {
			DefaultLogger.SetLevel(l)
		}
		if os.Getenv("LOG_FORMAT") == "json" {
			DefaultLogger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat:  "",
				DisableTimestamp: false,
				DataKey:          "",
				FieldMap:         nil,
				CallerPrettyfier: nil,
				PrettyPrint:      false,
			})
		}
		DefaultLogger.SetOutput(os.Stdout)
		DefaultBaseEntry = DefaultLogger.WithField("type", name)
	})
}

// LogError write log errors with message
func LogError(log *logrus.Entry, err error, message string) {
	log.WithError(err).Error("*** " + message + " ***")
}

// Tag sets a tag name then returns a log entry ready to write
func Tag(tag string) *logrus.Entry {
	if DefaultBaseEntry == nil {
		Init("function")
	}
	return DefaultBaseEntry.WithField("tag", tag)
}

// WithTag return a log entry from tag name
func WithTag(tag string) *logrus.Entry {
	l := Tag(tag)
	return l
}
