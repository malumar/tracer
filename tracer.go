package tracer

import (
	"fmt"
	"io"
	"os"
)

type Level int

const (
	Trace Level = 1 << iota
	Info
	Debug
	Warn
	Error
)

const All Level = Trace | Info | Debug | Warn | Error

type Tracer interface {
	NeedNoErr(err error) string
	SetValue(name string, value string)
	Write(b []byte) (n int, err error)
	WriteString(s string)
	Writef(s string, args ...interface{})
	WriteLine(s string)
	WriteLinef(s string, args ...interface{})
	Trace(message string, args ...interface{})
	Debug(message string, args ...interface{})
	Info(message string, args ...interface{})
	Error(message string, args ...interface{})
	Warn(message string, args ...interface{})
	Fatal(message string, args ...interface{})
	Close() error
}

func NewWithoutOnCloser(level Level, writer io.Writer) Tracer {
	return New(level, NopWriter{writer})
}

// New create new instance of tracer
// @level specify the warning levels you want to save, e.g. All for everyone or selected: Info|Error
// @writer logs will be saved here
func New(level Level, writer Writer) Tracer {

	return &tracer{
		level:  level,
		writer: writer,
	}
}

// NewSimple e.g. for using simple writing
// example using with slog:
//
//	NewSimple(All, func(bytes []byte) {
//		slog.Default().Info(string(bytes))
//	})
func NewSimple(level Level, handler func([]byte)) Tracer {
	return &tracer{
		level:  level,
		writer: NopWriter{handlerWriter{handler}},
	}
}

type handlerWriter struct {
	handler func([]byte)
}

func (self handlerWriter) Write(p []byte) (int, error) {
	self.handler(p)
	return len(p), nil
}

type NopWriter struct {
	io.Writer
}

func (NopWriter) OnClose() error { return nil }

type Writer interface {
	Write(b []byte) (n int, err error)
	OnClose() error
}

// Tracking performed procedures
type tracer struct {
	level  Level
	writer Writer
	values map[string]string
}

func (t *tracer) NeedNoErr(err error) string {
	if err == nil {
		return "[OK]"
	} else {
		return "[Error: " + err.Error() + "]"
	}
}
func (t *tracer) SetValue(name string, value string) {
	t.values[name] = value
}

func (t *tracer) Write(r []byte) (n int, err error) {
	return t.writer.Write(r)
}

func (t *tracer) WriteString(s string) {
	// Ignore unhandled error
	_, _ = t.Write([]byte(s))
}

func (t *tracer) Writef(s string, args ...interface{}) {
	t.WriteString(fmt.Sprintf(s, args...))
}

func (t *tracer) WriteLine(s string) {
	t.WriteString(s + "\n")
}

func (t *tracer) WriteLinef(s string, args ...interface{}) {
	t.WriteLine(fmt.Sprintf(s, args...))
}

func (t *tracer) Trace(message string, args ...interface{}) {
	if t.level&Trace != 0 {
		t.Writef(fmt.Sprintf("TRACE: %s\n", message), args...)
	}
}

func (t *tracer) Debug(message string, args ...interface{}) {
	if t.level&Debug != 0 {
		t.Writef(fmt.Sprintf("DEBUG: %s\n", message), args...)
	}
}

func (t *tracer) Info(message string, args ...interface{}) {
	if t.level&Info != 0 {
		t.Writef(fmt.Sprintf("INFO: %s\n", message), args...)

	}
}

func (t *tracer) Warn(message string, args ...interface{}) {
	if t.level&Warn != 0 {
		t.Writef(fmt.Sprintf("WARN: %s\n", message), args...)
	}
}

func (t *tracer) Error(message string, args ...interface{}) {
	if t.level&Error != 0 {
		t.Writef(fmt.Sprintf("ERROR: %s\n", message), args...)
	}
}

func (t *tracer) Fatal(message string, args ...interface{}) {
	t.Writef(fmt.Sprintf("FATAL: %s\n", message), args...)
	os.Exit(1)
}

func (t *tracer) Close() error {
	t.Debug("End of session")
	return t.writer.OnClose()
}
