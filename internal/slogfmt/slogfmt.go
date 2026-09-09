// This is for customizing how slog (structured log) will format the output in the terminal

package slogfmt

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func Init() {
	slog.SetDefault(slog.New(New(os.Stdout)))
}

type handler struct {
	w io.Writer
}

// New returns a slog.Handler that prints:
// [LEVEL] <time> <message> <file:line> key=value ...
func New(w io.Writer) slog.Handler {
	return &handler{w: w}
}

func (h *handler) Enabled(context.Context, slog.Level) bool { return true }

func (h *handler) Handle(_ context.Context, r slog.Record) error {
	var b strings.Builder
	b.WriteByte('[')
	b.WriteString(r.Level.String())
	b.WriteString("] ")
	b.WriteString(r.Time.Format("2006-01-02T15:04:05.000-07:00"))
	b.WriteByte(' ')
	b.WriteByte('"')
	b.WriteString(r.Message)
	b.WriteByte('"')

	if r.PC != 0 {
		f, _ := runtime.CallersFrames([]uintptr{r.PC}).Next()
		b.WriteByte(' ')
		b.WriteString(f.File)
		b.WriteByte(':')
		b.WriteString(strconv.Itoa(f.Line))
	}

	r.Attrs(func(a slog.Attr) bool {
		fmt.Fprintf(&b, " %s=%s", a.Key, a.Value.String())
		return true
	})

	_, _ = h.w.Write([]byte(b.String() + "\n"))
	return nil
}

func (h *handler) WithAttrs([]slog.Attr) slog.Handler { return h }
func (h *handler) WithGroup(string) slog.Handler      { return h }
