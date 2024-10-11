package tracer

import (
	"strings"
	"testing"
)

func TestNewSimple(t *testing.T) {

	var sb strings.Builder

	tr := NewSimple(Debug|Trace, func(bytes []byte) {
		sb.Write(bytes)
	})

	tr.Trace("Hello World")
	tr.Info("Hello world")
	tr.Error("Hello World")
	tr.Debug("Hello World")
	tr.Warn("Hello World")

	expected := "TRACE: Hello World\nDEBUG: Hello World\n"
	if sb.String() != expected {
		t.Errorf("no equal: expected:\n`%v` have:\n`%v`\n", expected, sb.String())
	}

	sb.Reset()
	tr = NewSimple(All, func(bytes []byte) {
		sb.Write(bytes)
	})

	tr.Trace("1")
	tr.Warn("2")
	tr.Info("3")
	tr.Debug("4")
	tr.Error("5")

	expected = "TRACE: 1\nWARN: 2\nINFO: 3\nDEBUG: 4\nERROR: 5\n"
	if sb.String() != expected {
		t.Errorf("no equal: expected:\n`%v` have:\n`%v`\n", expected, sb.String())
	}

}
