package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// ProcessLogger is a logger that is owned by a process.
type ProcessLogger struct {
	ProcessID  string
	Prefix     string
	stackTrace string
	entry      []*logrus.Entry
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = "36"
	gray    = 37
)

// NewProcessLogger instantiates a new process logger with a given process ID and log file path.
func NewProcessLogger(debug bool) (logger *ProcessLogger) {

	logger = &ProcessLogger{}
	logger.entry = []*logrus.Entry{}

	// Initialize logger.
	log := logrus.New()

	if debug {
		log.Level = logrus.DebugLevel
	}

	// Set color scheme.
	/* colors := &prefixed.ColorScheme{
		InfoLevelStyle: blue,
	} */

	// Initialize formatter.
	formatter := new(logrus.TextFormatter) // default
	formatter.DisableTimestamp = false
	formatter.FullTimestamp = false
	//formatter.SetColorScheme(colors)

	log.Formatter = formatter
	logger.entry = append(logger.entry, logrus.NewEntry(log))

	return
}

// AddJSONWriter adds a writer to the logger.
func (logger *ProcessLogger) AddJSONWriter(writer io.Writer) {

	log := logrus.New()
	log.Formatter = new(logrus.JSONFormatter)
	log.Out = writer
	logger.entry = append(logger.entry, logrus.NewEntry(log))

}

func (logger *ProcessLogger) writeFormattedLine(message string, callback func(*logrus.Entry, string)) {
	for i := range logger.entry {

		_, josnFormat := logger.entry[i].Logger.Formatter.(*logrus.JSONFormatter)

		if josnFormat {
			logger.entry[i].WithField("stack", logger.stackTrace)
		}

		if logger.Prefix != "" {
			logger.entry[i].WithField("prefix", logger.Prefix)
		}

		callback(logger.entry[i], message)

		if !josnFormat {
			io.WriteString(logger.entry[i].Logger.Out, logger.stackTrace)
		}
	}
}

// WriteDebug writes a debug message to the logger.
func (logger *ProcessLogger) WriteDebug(message string) {
	logger.writeFormattedLine(message, func(e *logrus.Entry, msg string) { e.Debugln(msg) })
}

// WriteInfo writes a debug message to the logger.
func (logger *ProcessLogger) WriteInfo(message string) {
	logger.writeFormattedLine(message, func(e *logrus.Entry, msg string) { e.Infoln(msg) })
}

// WriteWarning writes a debug message to the logger.
func (logger *ProcessLogger) WriteWarning(message string) {
	logger.writeFormattedLine(message, func(e *logrus.Entry, msg string) { e.Warnln(msg) })
}

// WriteError writes a debug message to the logger.
func (logger *ProcessLogger) WriteError(message string) {
	logger.writeFormattedLine(message, func(e *logrus.Entry, msg string) { e.Errorln(msg) })

}

// WriteFatal writes a debug message to the logger.
func (logger *ProcessLogger) WriteFatal(message string) {
	logger.writeFormattedLine(message, func(e *logrus.Entry, msg string) { e.Fatalln(msg) })
}

// DeepCopy deepcopies a to b using json marshaling
func SemiDeepCopy(a, b interface{}) {
	byt, _ := json.Marshal(a)
	json.Unmarshal(byt, b)
}

func CopyLogger(logger ProcessLogger) ProcessLogger{
	//SemiDeepCopy does not follow pointers or copy map elements
	//Because of this we could manually copy the two fields but this is more resilient to changes
	var result ProcessLogger
	SemiDeepCopy(logger,result)
	result.entry = make([]*logrus.Entry,len(logger.entry))

	//Abuse of WithFields(nil) to make a copy of the entries
	//They reference to the same Logger within the entry
	for i := range logger.entry {
		result.entry[i] = logger.entry[i].WithFields(nil)
	}
	return result
}

// WithFields adds fields to the next logged message.
func (logger *ProcessLogger) WithFields(args ...interface{}) Logger {

	fields := logrus.Fields{}

	for i := 0; i < len(args); i += 2 {
		fields[args[i].(string)] = args[i+1]
	}

	//Lighter than result:=CopyLogger(*logger)
	var result ProcessLogger
	SemiDeepCopy(logger,result)
	result.entry = make([]*logrus.Entry,len(logger.entry))

	for i := range logger.entry {
		result.entry[i] = logger.entry[i].WithFields(fields)
	}

	return &result
}

// WithStack adds a stack trace from a given error.
func (logger *ProcessLogger) WithStack(err error) Logger {

	result:=CopyLogger(*logger)
	
	if err != nil {

		var builder strings.Builder

		if err, ok := err.(stackTracer); ok {
			builder.WriteString("\nStack Trace:\n")
			for _, f := range err.StackTrace() {
				builder.WriteString(fmt.Sprintf("%+v\n", f))
			}
			builder.WriteString("\n")
		}

		result.stackTrace += builder.String()
	}

	return &result
}

// WithError adds an error message from a given error.
func (logger *ProcessLogger) WithError(err error) Logger {

	//Lighter than result:=CopyLogger(*logger)
	var result ProcessLogger
	SemiDeepCopy(logger,result)
	result.entry = make([]*logrus.Entry,len(logger.entry))

	if err != nil {
		for i := range logger.entry {
			result.entry[i] = (*logger.entry[i]).WithField("error", err.Error())
		}
	}
	return &result
}
