package slogdefault

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"sync"
)

type defaultHandler struct {
	mu   sync.Mutex
	next slog.Handler
	buf  *bytes.Buffer
	w    io.Writer
}

// DefaultReplaceAttr makes your logs look how slog.Default logs look
func ReplaceAttr(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey {
		return slog.Attr{}
	}
	if a.Key == slog.MessageKey {
		return slog.Attr{}
	}
	if a.Key == slog.LevelKey {
		return slog.Attr{}
	}
	return a
}

func NewHandler(w io.Writer, opts *slog.HandlerOptions) slog.Handler {
	dh := &defaultHandler{
		buf: new(bytes.Buffer),
		w:   w,
	}
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	if opts.ReplaceAttr == nil {
		opts.ReplaceAttr = ReplaceAttr
	}
	next := slog.NewTextHandler(dh.buf, opts)
	dh.next = next
	return dh
}

func (d *defaultHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return d.next.Enabled(ctx, l)
}

// Collect the level, attributes and message in a string and
// write it with the default log.Logger.
// Let the log.Logger handle time and file/line.
func (h *defaultHandler) Handle(ctx context.Context, r slog.Record) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	buf := new(bytes.Buffer)
	buf.WriteString(r.Time.Format("2006/01/02 15:04:05"))
	buf.WriteByte(' ')
	buf.WriteString(r.Level.String())
	buf.WriteByte(' ')
	buf.WriteString(r.Message)
	buf.WriteByte(' ')
	// combined with the ReplaceAttr above which strips time, message, and level
	// (handled above) - we just need TextHandler to output the attrs
	if err := h.next.Handle(ctx, r); err != nil {
		return err
	}
	hbufbytes := h.buf.Bytes()
	buf.Write(hbufbytes)
	bufbytes := buf.Bytes()
	_, err := h.w.Write(bufbytes)
	h.buf.Reset()
	return err
}

func (h *defaultHandler) WithAttrs(as []slog.Attr) slog.Handler {
	return &defaultHandler{w: h.w, buf: h.buf, next: h.next.WithAttrs(as)}
}

func (h *defaultHandler) WithGroup(name string) slog.Handler {
	return &defaultHandler{w: h.w, buf: h.buf, next: h.next.WithGroup(name)}
}
