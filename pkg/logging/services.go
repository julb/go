package logging

import (
	"io"
	"log/syslog"
	"os"
	"time"

	logrus "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
)

var (
	globalFields           map[string]interface{}
	standardLogWithContext *LogWithContext
)

type LogWithContext struct {
	logrusEntry *logrus.Entry
}

type LoggingOpts struct {
	Level string `yaml:"level" json:"level"`

	Output struct {
		Console struct {
			Enabled bool   `yaml:"enabled" json:"enabled"`
			Type    string `yaml:"type" json:"type"`
		} `yaml:"console" json:"console"`

		File struct {
			Enabled bool   `yaml:"enabled" json:"enabled"`
			Path    string `yaml:"path" json:"path"`
		} `yaml:"file" json:"file"`

		Syslog struct {
			Enabled  bool            `yaml:"enabled" json:"enabled"`
			Protocol string          `yaml:"protocol" json:"protocol"`
			Addr     string          `yaml:"addr" json:"addr"`
			Priority syslog.Priority `yaml:"priority" json:"priority"`
			Tag      string          `yaml:"tag" json:"tag"`
		} `yaml:"syslog" json:"syslog"`
	} `yaml:"output" json:"output"`

	Formatter struct {
		TextFormatter struct {
			Enabled bool `yaml:"enabled" json:"enabled"`
		} `yaml:"text" json:"text"`

		JSONFormatter struct {
			Enabled      bool                   `yaml:"enabled" json:"enabled"`
			GlobalFields map[string]interface{} `yaml:"globalFields" json:"globalFields"`
		} `yaml:"json" json:"json"`
	} `yaml:"formatter" json:"formatter"`
}

func init() {
	// Initialize default logger.
	standardLogWithContext = WithEmptyContext()

	// Configure the logging context with defaults options.
	Configure(DefaultLoggingOpts()) // nolint:errcheck
}

// The default logging opts.
func DefaultLoggingOpts() *LoggingOpts {
	defaultLoggingOpts := &LoggingOpts{
		Level: "info",
	}

	defaultLoggingOpts.Output.Console.Enabled = true
	defaultLoggingOpts.Output.Console.Type = "stdout"
	defaultLoggingOpts.Formatter.TextFormatter.Enabled = true

	return defaultLoggingOpts
}

// Configure logging context
func Configure(opts *LoggingOpts) error {
	SetLevel(opts.Level)

	// Console
	var consoleOutput io.Writer
	if opts.Output.Console.Enabled {
		if opts.Output.Console.Type == "stdout" {
			consoleOutput = os.Stdout
		} else if opts.Output.Console.Type == "stderr" {
			consoleOutput = os.Stderr
		}
	}

	// File
	var fileOutput io.Writer
	var err error
	if opts.Output.File.Enabled {
		fileOutput, err = os.OpenFile(opts.Output.File.Path, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return err
		}
	}

	if consoleOutput != nil && fileOutput != nil {
		logrus.SetOutput(io.MultiWriter(consoleOutput, fileOutput))
	} else if consoleOutput != nil {
		logrus.SetOutput(consoleOutput)
	} else if fileOutput != nil {
		logrus.SetOutput(fileOutput)
	}

	// Syslog output.
	if opts.Output.Syslog.Enabled {
		err = AddSyslogOutput(opts.Output.Syslog.Protocol, opts.Output.Syslog.Addr, opts.Output.Syslog.Priority, opts.Output.Syslog.Tag)
		if err != nil {
			return err
		}
	}

	if opts.Formatter.TextFormatter.Enabled {
		SetTextFormatter()
	} else if opts.Formatter.JSONFormatter.Enabled {
		SetJSONFormatter(DefaultJSONFormatterOpts())
		SetGlobalFields(opts.Formatter.JSONFormatter.GlobalFields)
	}

	return nil
}

// Set the current level of logging
func SetLevel(level string) {
	// Get logrus level
	logrusLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logrus.Fatal(err.Error())
	}

	// Set level in logrus.
	logrus.SetLevel(logrusLevel)
}

// Panic logs a message at level Panic on the standard logger.
func Panic(args ...interface{}) {
	standardLogWithContext.Panic(args...)
}

// Panicf logs a message at level Panic on the standard logger.
func Panicf(format string, args ...interface{}) {
	standardLogWithContext.Panicf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatal(args ...interface{}) {
	standardLogWithContext.Fatal(args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func Fatalf(format string, args ...interface{}) {
	standardLogWithContext.Fatalf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func Error(args ...interface{}) {
	standardLogWithContext.Error(args...)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	standardLogWithContext.Errorf(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func Warn(args ...interface{}) {
	standardLogWithContext.Warn(args...)
}

// Warf logs a message at level Warn on the standard logger.
func Warnf(format string, args ...interface{}) {
	standardLogWithContext.Warnf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func Info(args ...interface{}) {
	standardLogWithContext.Info(args...)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	standardLogWithContext.Infof(format, args...)
}

// Debug logs a message at level Debug on the standard logger.
func Debug(args ...interface{}) {
	standardLogWithContext.Debug(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	standardLogWithContext.Debugf(format, args...)
}

// Trace logs a message at level Trace on the standard logger.
func Trace(args ...interface{}) {
	standardLogWithContext.Trace(args...)
}

// Tracef logs a message at level Trace on the standard logger.
func Tracef(format string, args ...interface{}) {
	standardLogWithContext.Tracef(format, args...)
}

// Return a logger with specified field added.
func WithEmptyContext() *LogWithContext {
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	return &LogWithContext{
		logrusEntry: logrusEntry,
	}
}

// Return a logger with specified field added.
func WithField(key string, value interface{}) *LogWithContext {
	logrusEntry := logrus.WithFields(globalFields).WithField(key, value)
	return &LogWithContext{
		logrusEntry: logrusEntry,
	}
}

// Return a logger with request ID field present
func WithRequestId(value string) *LogWithContext {
	return WithField("x-request-id", value)
}

// Return a logger with trademark field present
func WithTrademark(value string) *LogWithContext {
	return WithField("tm", value)
}

// Return a logger with trace ID and span ID fields present
func WithTracing(traceId string, spanId string) *LogWithContext {
	return WithField("x-b3-traceid", traceId).WithField("x-b3-spanid", spanId)
}

// Panic logs a message at level Panic on the standard logger.
func (logWithContext *LogWithContext) Panic(args ...interface{}) {
	logWithContext.logrusEntry.Panic(args...)
}

// Panicf logs a message at level Panic on the standard logger.
func (logWithContext *LogWithContext) Panicf(format string, args ...interface{}) {
	logWithContext.logrusEntry.Panicf(format, args...)
}

// Fatal logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (logWithContext *LogWithContext) Fatal(args ...interface{}) {
	logWithContext.logrusEntry.Fatal(args...)
}

// Fatalf logs a message at level Fatal on the standard logger then the process will exit with status set to 1.
func (logWithContext *LogWithContext) Fatalf(format string, args ...interface{}) {
	logWithContext.logrusEntry.Fatalf(format, args...)
}

// Error logs a message at level Error on the standard logger.
func (logWithContext *LogWithContext) Error(args ...interface{}) {
	logWithContext.logrusEntry.Error(args...)
}

// Errorf logs a message at level Error on the standard logger.
func (logWithContext *LogWithContext) Errorf(format string, args ...interface{}) {
	logWithContext.logrusEntry.Errorf(format, args...)
}

// Warn logs a message at level Warn on the standard logger.
func (logWithContext *LogWithContext) Warn(args ...interface{}) {
	logWithContext.logrusEntry.Warn(args...)
}

// Warf logs a message at level Warn on the standard logger.
func (logWithContext *LogWithContext) Warnf(format string, args ...interface{}) {
	logWithContext.logrusEntry.Warnf(format, args...)
}

// Info logs a message at level Info on the standard logger.
func (logWithContext *LogWithContext) Info(args ...interface{}) {
	logWithContext.logrusEntry.Info(args...)
}

// Infof logs a message at level Info on the standard logger.
func (logWithContext *LogWithContext) Infof(format string, args ...interface{}) {
	logWithContext.logrusEntry.Infof(format, args...)
}

// Debug logs a message at level Debug on the standard logger.
func (logWithContext *LogWithContext) Debug(args ...interface{}) {
	logWithContext.logrusEntry.Debug(args...)
}

// Debugf logs a message at level Debug on the standard logger.
func (logWithContext *LogWithContext) Debugf(format string, args ...interface{}) {
	logWithContext.logrusEntry.Debugf(format, args...)
}

// Trace logs a message at level Trace on the standard logger.
func (logWithContext *LogWithContext) Trace(args ...interface{}) {
	logWithContext.logrusEntry.Trace(args...)
}

// Tracef logs a message at level Trace on the standard logger.
func (logWithContext *LogWithContext) Tracef(format string, args ...interface{}) {
	logWithContext.logrusEntry.Tracef(format, args...)
}

// Return a new logger with request ID field present
func (logWithContext *LogWithContext) WithRequestId(value string) *LogWithContext {
	return logWithContext.WithField("x-request-id", value)
}

// Return a new logger with trademark field present
func (logWithContext *LogWithContext) WithTrademark(value string) *LogWithContext {
	return logWithContext.WithField("tm", value)
}

// Return a logger with trace ID and span ID fields present
func (logWithContext *LogWithContext) WithTracing(traceId string, spanId string) *LogWithContext {
	return logWithContext.WithField("x-b3-traceid", traceId).WithField("x-b3-spanid", spanId)
}

// Return a new logger with specified field added.
func (logWithContext *LogWithContext) WithField(key string, value interface{}) *LogWithContext {
	newLogrusEntry := logWithContext.logrusEntry.WithField(key, value)
	return &LogWithContext{
		logrusEntry: newLogrusEntry,
	}
}

// Sets the log output to stdout.
func SetOutputToStdout() {
	SetOutput(os.Stdout)
}

// Set the log output to stderr.
func SetOutputToStderr() {
	SetOutput(os.Stderr)
}

// Set the log output to the given IO Writer
func SetOutput(out io.Writer) {
	logrus.SetOutput(out)
}

// Sets global fields to the logger.
// Clear any existing values and set the new values.
func SetGlobalFields(globalFieldsToSet map[string]interface{}) {
	globalFields = make(map[string]interface{})
	AddGlobalFields(globalFieldsToSet)
}

// Add global fields to the logger.
func AddGlobalFields(globalFieldsToAdd map[string]interface{}) {
	// Add fields to the map.
	for k, v := range globalFieldsToAdd {
		globalFields[k] = v
	}

	// Create a fresh new logrus entry.
	standardLogWithContext.logrusEntry = logrus.NewEntry(logrus.StandardLogger()).WithFields(globalFields)
}

// Log to syslog server
//
// - network and raddr: see https://pkg.go.dev/net#Dial function for a list of supported values.
//
// - priority: syslog priority for the message
//
// - tag: tag to assign to the log
func AddSyslogOutput(network string, raddr string, priority syslog.Priority, tag string) error {
	// Trace
	logrus.Debugf("adding syslog output with protocol=%s | addr=%s | priority=%d | tag=%s.", network, raddr, priority, tag)

	// Add syslog hook
	hook, err := logrus_syslog.NewSyslogHook(network, raddr, priority, tag)
	if err != nil {
		logrus.Errorf("unable to connect to syslog %s/%s.", network, raddr)
		return err
	}

	// Add hook
	logrus.AddHook(hook)

	// Trace successful.
	logrus.Debugf("syslog output with protocol=%s | addr=%s | priority=%d | tag=%s configured successfully", network, raddr, priority, tag)
	return nil
}

// Sets a format of output to Text
func SetTextFormatter() {
	logrus.SetFormatter(&logrus.TextFormatter{})
}

// JSONFormatter formats logs into parsable json
type JSONFormatterOpts struct {
	// TimestampFormat sets the format used for marshaling timestamps.
	// The format to use is the same than for time.Format or time.Parse from the standard
	// library.
	// The standard Library already provides a set of predefined format.
	TimestampFormat string

	// DataKey allows users to put all the log entry parameters into a nested dictionary at a given key.
	DataKey string

	// TimeKey allows users to set a specific JSON key for the time.
	TimestampKey string

	// LevelKey allows users to set a specific JSON key for the level.
	LevelKey string

	// MessageKey allows users to set a specific JSON key for the message.
	MessageKey string
}

// Get a set of JSON formatter options by default.
func DefaultJSONFormatterOpts() *JSONFormatterOpts {
	return &JSONFormatterOpts{
		TimestampFormat: time.RFC3339Nano,
		TimestampKey:    "time",
		LevelKey:        "level",
		MessageKey:      "message",
	}
}

// Sets a formatter of output to JSON
func SetJSONFormatter(o *JSONFormatterOpts) {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: o.TimestampFormat,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  o.TimestampKey,
			logrus.FieldKeyLevel: o.LevelKey,
			logrus.FieldKeyMsg:   o.MessageKey,
		},
		DataKey: o.DataKey,
	})
}
