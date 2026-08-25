package zerolog

import (
	"io"
)

type Logger struct {
	w       io.Writer
	level   Level
	context []byte
	hooks   []Hook
}

type Hook interface {
	Run(e *Event, level Level, message string)
}

func (l Logger) Output(w io.Writer) Logger {
	l.w = w
	return l
}

func (l Logger) Level(lvl Level) Logger {
	l.level = lvl
	return l
}

func (l Logger) GetLevel() Level {
	return l.level
}

func (l Logger) Hook(h Hook) Logger {
	l.hooks = append(l.hooks[:len(l.hooks):len(l.hooks)], h)
	return l
}

// With creates a child logger context with capacity bounded to current length.
func (l Logger) With() Context {
	context := l.context
	if context != nil {
		context = context[:len(context):len(context)]
	}
	return Context{l: l, context: context}
}

func (l Logger) newEvent(level Level) *Event {
	if level < l.level || l.level == Disabled {
		return nil
	}
	e := newEvent(l.w, level)
	if len(l.context) > 0 {
		ctx := l.context
		if len(e.buf) == 1 { // buffer currently contains only '{'
			if ctx[0] == ',' {
				ctx = ctx[1:]
			}
		} else {
			if ctx[0] != ',' {
				e.buf = append(e.buf, ',')
			}
		}
		e.buf = append(e.buf, ctx...)
	}
	e.hooks = l.hooks
	return e
}

func (l Logger) Trace() *Event {
	return l.newEvent(TraceLevel)
}

func (l Logger) Debug() *Event {
	return l.newEvent(DebugLevel)
}

func (l Logger) Info() *Event {
	return l.newEvent(InfoLevel)
}

func (l Logger) Warn() *Event {
	return l.newEvent(WarnLevel)
}

func (l Logger) Error() *Event {
	return l.newEvent(ErrorLevel)
}

func (l Logger) Fatal() *Event {
	return l.newEvent(FatalLevel)
}

func (l Logger) Panic() *Event {
	return l.newEvent(PanicLevel)
}

func (l Logger) Log() *Event {
	return l.newEvent(NoLevel)
}

func (l Logger) Write(p []byte) (n int, err error) {
	if l.w != nil {
		return l.w.Write(p)
	}
	return len(p), nil
}
