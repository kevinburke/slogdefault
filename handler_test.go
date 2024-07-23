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

func TestWith(t *testing.T) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	buf := new(bytes.Buffer)
	ch := NewHandler(buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	l := slog.New(ch)
	log_(l)
	l2 := l.With("blah", 26)
	log_(l2)
	if count := bytes.Count(buf.Bytes(), []byte{'\n'}); count != 8 {
		t.Errorf("expected 8 newlines in log output, got %v", count)
	}
	if count := strings.Count(buf.String(), "blah=26"); count != 4 {
		t.Errorf("expected sub attribute to be present 4 times, was not: %d", count)
	}
	if !strings.Contains(buf.String(), "blah=26") {
		t.Errorf("expected sub attribute to be present, was not")
	}
}

func TestDefaultTextHandlerWith(t *testing.T) {
	buf := new(bytes.Buffer)
	next := slog.NewTextHandler(buf, nil)
	l := slog.New(next)
	log_(l)
	l2 := l.With("blah", 37)
	log_(l2)
	if count := bytes.Count(buf.Bytes(), []byte{'\n'}); count != 8 {
		t.Errorf("expected 8 newlines in log output, got %v", count)
	}
}
