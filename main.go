package zerolog

import (
	"io"
	"os"
	"time"
)

type Level int8

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
	PanicLevel
	NoLevel
	Disabled
	TraceLevel Level = -1
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case Disabled:
		return "disabled"
	default:
		return ""
	}
}

var (
	TimestampFunc      = time.Now
	TimestampFieldName = "time"
	LevelFieldName     = "level"
	MessageFieldName   = "message"
	ErrorFieldName     = "error"
	TimeFieldFormat    = time.RFC3339

	DefaultContextLogger *Logger
)

func init() {
	l := New(os.Stderr)
	DefaultContextLogger = &l
}

func New(w io.Writer) Logger {
	if w == nil {
		w = io.Discard
	}
	return Logger{w: w, level: TraceLevel}
}

func Nop() Logger {
	return Logger{w: io.Discard, level: Disabled}
}
