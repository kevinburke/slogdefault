# slogdefault

This is a trivial library that exposes the standard library's
`log/slog.Default()` and makes it customizable.

### Usage

```go
h := slogdefault.NewHandler(output, &slog.HandlerOptions{/* ... */}
logger := slog.New(h)
logger.Info("test message")
```

Use this to add source lines, change the log level, chain handlers, etc.
