package zerolog

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
	"time"
)

var eventPool = sync.Pool{
	New: func() interface{} {
		return &Event{
			buf: make([]byte, 0, 512),
		}
	},
}

type Event struct {
	w     io.Writer
	level Level
	buf   []byte
	hooks []Hook
}

func newEvent(w io.Writer, level Level) *Event {
	e := eventPool.Get().(*Event)
	e.w = w
	e.level = level
	e.buf = e.buf[:0]
	e.hooks = nil

	e.buf = append(e.buf, '{')
	if level != NoLevel {
		e.buf = encAppendKeyNoComma(e.buf, LevelFieldName)
		e.buf = encAppendString(e.buf, level.String())
	}
	return e
}

func encAppendKeyNoComma(dst []byte, key string) []byte {
	dst = encAppendString(dst, key)
	dst = append(dst, ':')
	return dst
}

func (e *Event) Enabled() bool {
	return e != nil
}

func (e *Event) Str(key, val string) *Event {
	if e == nil {
		return nil
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = encAppendString(e.buf, val)
	return e
}

func (e *Event) Bytes(key string, val []byte) *Event {
	if e == nil {
		return nil
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = encAppendBytes(e.buf, val)
	return e
}

func (e *Event) Hex(key string, val []byte) *Event {
	if e == nil {
		return nil
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = encAppendHex(e.buf, val)
	return e
}

func (e *Event) RawJSON(key string, b []byte) *Event {
	if e == nil {
		return nil
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = append(e.buf, b...)
	return e
}

func (e *Event) Int(key string, val int) *Event {
	return e.Int64(key, int64(val))
}

func (e *Event) Int8(key string, val int8) *Event {
	return e.Int64(key, int64(val))
}

func (e *Event) Int16(key string, val int16) *Event {
	return e.Int64(key, int64(val))
}

func (e *Event) Int32(key string, val int32) *Event {
	return e.Int64(key, int64(val))
}

func (e *Event) Int64(key string, val int64) *Event {
	if e == nil {
		return nil
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = strconv.AppendInt(e.buf, val, 10)
	return e
}

func (e *Event) Uint(key string, val uint) *Event {
	return e.Uint64(key, uint64(val))
}

func (e *Event) Uint8(key string, val uint8) *Event {
	return e.Uint64(key, uint64(val))
}

func (e *Event) Uint16(key string, val uint16) *Event {
	return e.Uint64(key, uint64(val))
}

func (e *Event) Uint32(key string, val uint32) *Event {
	return e.Uint64(key, uint64(val))
}

func (e *Event) Uint64(key string, val uint64) *Event {
	if e == nil {
		return nil
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = strconv.AppendUint(e.buf, val, 10)
	return e
}

func (e *Event) Float32(key string, val float32) *Event {
	return e.Float64(key, float64(val))
}

func (e *Event) Float64(key string, val float64) *Event {
	if e == nil {
		return nil
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = strconv.AppendFloat(e.buf, val, 'f', -1, 64)
	return e
}

func (e *Event) Bool(key string, val bool) *Event {
	if e == nil {
		return nil
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = strconv.AppendBool(e.buf, val)
	return e
}

func (e *Event) Err(err error) *Event {
	if e == nil || err == nil {
		return e
	}
	return e.Str(ErrorFieldName, err.Error())
}

func (e *Event) AnErr(key string, err error) *Event {
	if e == nil || err == nil {
		return e
	}
	return e.Str(key, err.Error())
}

func (e *Event) Timestamp() *Event {
	if e == nil {
		return nil
	}
	t := TimestampFunc()
	if TimeFieldFormat == "" {
		return e.Int64(TimestampFieldName, t.Unix())
	}
	return e.Str(TimestampFieldName, t.Format(TimeFieldFormat))
}

func (e *Event) Time(key string, t time.Time) *Event {
	if e == nil {
		return nil
	}
	if TimeFieldFormat == "" {
		return e.Int64(key, t.Unix())
	}
	return e.Str(key, t.Format(TimeFieldFormat))
}

func (e *Event) Dur(key string, d time.Duration) *Event {
	if e == nil {
		return nil
	}
	return e.Int64(key, d.Nanoseconds())
}

func (e *Event) Interface(key string, val interface{}) *Event {
	if e == nil {
		return nil
	}
	b, err := json.Marshal(val)
	if err != nil {
		return e.Str(key, fmt.Sprintf("[json error: %v]", err))
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = append(e.buf, b...)
	return e
}

func (e *Event) Any(key string, val interface{}) *Event {
	return e.Interface(key, val)
}

func (e *Event) Fields(fields map[string]interface{}) *Event {
	if e == nil {
		return nil
	}
	for k, v := range fields {
		e.Interface(k, v)
	}
	return e
}

func (e *Event) Dict(key string, dict Context) *Event {
	if e == nil {
		return nil
	}
	if len(e.buf) > 1 {
		e.buf = append(e.buf, ',')
	}
	e.buf = encAppendKeyNoComma(e.buf, key)
	e.buf = append(e.buf, '{')
	if len(dict.context) > 0 {
		d := dict.context
		if d[0] == ',' {
			d = d[1:]
		}
		e.buf = append(e.buf, d...)
	}
	e.buf = append(e.buf, '}')
	return e
}

func (e *Event) Msg(msg string) {
	if e == nil {
		return
	}
	e.msg(msg)
}

func (e *Event) Msgf(format string, v ...interface{}) {
	if e == nil {
		return
	}
	e.msg(fmt.Sprintf(format, v...))
}

func (e *Event) Send() {
	if e == nil {
		return
	}
	e.msg("")
}

func (e *Event) msg(msg string) {
	for _, hook := range e.hooks {
		hook.Run(e, e.level, msg)
	}
	if msg != "" {
		if len(e.buf) > 1 {
			e.buf = append(e.buf, ',')
		}
		e.buf = encAppendKeyNoComma(e.buf, MessageFieldName)
		e.buf = encAppendString(e.buf, msg)
	}
	e.buf = append(e.buf, '}', '\n')
	if e.w != nil {
		_, _ = e.w.Write(e.buf)
	}
	lvl := e.level
	if cap(e.buf) <= 4096 {
		eventPool.Put(e)
	}
	if lvl == FatalLevel {
		os.Exit(1)
	} else if lvl == PanicLevel {
		panic(msg)
	}
}
