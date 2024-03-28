package slogdefault

import (
	"bytes"
	"log"
	"log/slog"
	"strings"
	"testing"
)

func log_(l *slog.Logger) {
	l.Info("test message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
	l.Warn("warn message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
	l.Error("error message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
	l.Debug("debug message", "foo", "bar", "one", 2, "two", []int{3, 4, 5})
}

func TestDefaultHandler(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	buf := new(bytes.Buffer)
	ch := NewHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := slog.New(ch)
	log_(l)
	buf2 := new(bytes.Buffer)
	log.SetOutput(buf2)
	log_(slog.Default())
	if !bytes.Equal(buf.Bytes(), buf2.Bytes()) {
		t.Errorf("expected two buffers to have same output:\n\n%s\n%s", buf.String(), buf2.String())
	}
}

func TestAddSource(t *testing.T) {
	buf := new(bytes.Buffer)
	ch := NewHandler(buf, &slog.HandlerOptions{
		Level:     slog.LevelWarn,
		AddSource: true,
	})
	log_(slog.New(ch))
	if !strings.Contains(buf.String(), `slogdefault/handler_test.go:14`) {
		t.Errorf("source log did not match expected string %q", buf.String())
	}
}
