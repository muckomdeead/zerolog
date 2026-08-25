package zerolog

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Context struct {
	l       Logger
	context []byte
}

func (c Context) Logger() Logger {
	l := c.l
	if len(c.context) > 0 {
		l.context = c.context[:len(c.context):len(c.context)]
	} else {
		l.context = nil
	}
	return l
}

func (c Context) Str(key, val string) Context {
	c.context = encAppendKey(c.context, key)
	c.context = encAppendString(c.context, val)
	return c
}

func (c Context) Bytes(key string, val []byte) Context {
	c.context = encAppendKey(c.context, key)
	c.context = encAppendBytes(c.context, val)
	return c
}

func (c Context) Hex(key string, val []byte) Context {
	c.context = encAppendKey(c.context, key)
	c.context = encAppendHex(c.context, val)
	return c
}

func (c Context) RawJSON(key string, b []byte) Context {
	c.context = encAppendKey(c.context, key)
	c.context = append(c.context, b...)
	return c
}

func (c Context) Int(key string, val int) Context {
	return c.Int64(key, int64(val))
}

func (c Context) Int8(key string, val int8) Context {
	return c.Int64(key, int64(val))
}

func (c Context) Int16(key string, val int16) Context {
	return c.Int64(key, int64(val))
}

func (c Context) Int32(key string, val int32) Context {
	return c.Int64(key, int64(val))
}

func (c Context) Int64(key string, val int64) Context {
	c.context = encAppendKey(c.context, key)
	c.context = strconv.AppendInt(c.context, val, 10)
	return c
}

func (c Context) Uint(key string, val uint) Context {
	return c.Uint64(key, uint64(val))
}

func (c Context) Uint8(key string, val uint8) Context {
	return c.Uint64(key, uint64(val))
}

func (c Context) Uint16(key string, val uint16) Context {
	return c.Uint64(key, uint64(val))
}

func (c Context) Uint32(key string, val uint32) Context {
	return c.Uint64(key, uint64(val))
}

func (c Context) Uint64(key string, val uint64) Context {
	c.context = encAppendKey(c.context, key)
	c.context = strconv.AppendUint(c.context, val, 10)
	return c
}

func (c Context) Float32(key string, val float32) Context {
	return c.Float64(key, float64(val))
}

func (c Context) Float64(key string, val float64) Context {
	c.context = encAppendKey(c.context, key)
	c.context = strconv.AppendFloat(c.context, val, 'f', -1, 64)
	return c
}

func (c Context) Bool(key string, val bool) Context {
	c.context = encAppendKey(c.context, key)
	c.context = strconv.AppendBool(c.context, val)
	return c
}

func (c Context) Err(err error) Context {
	if err == nil {
		return c
	}
	return c.Str(ErrorFieldName, err.Error())
}

func (c Context) AnErr(key string, err error) Context {
	if err == nil {
		return c
	}
	return c.Str(key, err.Error())
}

func (c Context) Timestamp() Context {
	c.context = encAppendKey(c.context, TimestampFieldName)
	t := TimestampFunc()
	if TimeFieldFormat == "" {
		c.context = strconv.AppendInt(c.context, t.Unix(), 10)
	} else {
		c.context = encAppendString(c.context, t.Format(TimeFieldFormat))
	}
	return c
}

func (c Context) Time(key string, t time.Time) Context {
	c.context = encAppendKey(c.context, key)
	if TimeFieldFormat == "" {
		c.context = strconv.AppendInt(c.context, t.Unix(), 10)
	} else {
		c.context = encAppendString(c.context, t.Format(TimeFieldFormat))
	}
	return c
}

func (c Context) Dur(key string, d time.Duration) Context {
	c.context = encAppendKey(c.context, key)
	c.context = strconv.AppendInt(c.context, d.Nanoseconds(), 10)
	return c
}

func (c Context) Interface(key string, val interface{}) Context {
	c.context = encAppendKey(c.context, key)
	b, err := json.Marshal(val)
	if err != nil {
		c.context = encAppendString(c.context, fmt.Sprintf("[json error: %v]", err))
	} else {
		c.context = append(c.context, b...)
	}
	return c
}

func (c Context) Any(key string, val interface{}) Context {
	return c.Interface(key, val)
}

func (c Context) Fields(fields map[string]interface{}) Context {
	for k, v := range fields {
		c = c.Interface(k, v)
	}
	return c
}

func (c Context) Dict(key string, dict Context) Context {
	c.context = encAppendKey(c.context, key)
	c.context = append(c.context, '{')
	if len(dict.context) > 0 {
		d := dict.context
		if d[0] == ',' {
			d = d[1:]
		}
		c.context = append(c.context, d...)
	}
	c.context = append(c.context, '}')
	return c
}

func encAppendKey(dst []byte, key string) []byte {
	dst = append(dst, ',')
	dst = encAppendString(dst, key)
	dst = append(dst, ':')
	return dst
}

func encAppendString(dst []byte, s string) []byte {
	b, _ := json.Marshal(s)
	return append(dst, b...)
}

func encAppendBytes(dst []byte, s []byte) []byte {
	b, _ := json.Marshal(string(s))
	return append(dst, b...)
}

func encAppendHex(dst []byte, s []byte) []byte {
	return encAppendString(dst, fmt.Sprintf("%x", s))
}
