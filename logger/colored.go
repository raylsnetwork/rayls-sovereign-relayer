package logger

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"sync"

	"github.com/fatih/color"
)

var faintWhite = color.New(color.FgWhite, color.Faint)

type Options struct {
	Level slog.Leveler
}

type ColoredHandler struct {
	opts               Options
	preformattedGroups string
	preformattedAttrs  []byte

	projectRoot string

	mu  *sync.Mutex
	out io.Writer
}

func NewColored(out io.Writer, opts *Options) *ColoredHandler {
	projectRoot, err := getProjectRoot()
	if err != nil {
		projectRoot = "/"
	}

	h := &ColoredHandler{
		out: out,
		mu:  new(sync.Mutex),

		projectRoot: projectRoot,
	}

	if opts != nil {
		h.opts = *opts
	}
	if h.opts.Level == nil {
		h.opts.Level = slog.LevelInfo
	}
	return h
}

func (h *ColoredHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *ColoredHandler) Handle(ctx context.Context, r slog.Record) error {
	buf := make([]byte, 0, 1024)
	if !r.Time.IsZero() {
		buf = r.Time.AppendFormat(buf, "15:04:05 2006-01-02")
		buf = append(buf, ' ')
	}
	if r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			src := slog.Source{
				Function: f.Function,
				File:     f.File,
				Line:     f.Line,
			}
			buf = h.appendSource(buf, src)
			buf = append(buf, ' ')
		}
	}

	buf = h.appendLevel(buf, r.Level)
	buf = append(buf, ' ')

	buf = append(buf, r.Message...)
	buf = append(buf, " | "...)

	if len(h.preformattedAttrs) != 0 {
		buf = append(buf, h.preformattedAttrs...)
		buf = append(buf, ' ')
	}

	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, a, h.preformattedGroups)
		buf = append(buf, ' ')
		return true
	})

	buf = append(buf, '\n')

	h.mu.Lock()
	defer h.mu.Unlock()
	_, err := h.out.Write(buf)
	if err != nil {
		return fmt.Errorf("writing log output: %w", err)
	}
	return nil
}

func (h *ColoredHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	h2 := *h

	pre := fmt.Sprintf("%s%s.", h.preformattedGroups, name)
	h2.preformattedGroups = pre

	return &h2
}

func (h *ColoredHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	h2 := *h

	pre := slices.Clip(h.preformattedAttrs)
	for _, attr := range attrs {
		h2.preformattedAttrs = h2.appendAttr(pre, attr, h.preformattedGroups)
	}

	return &h2
}

func (h *ColoredHandler) appendAttr(buf []byte, a slog.Attr, prefix string) []byte {
	if a.Equal(slog.Attr{}) {
		return buf
	}
	if len(buf) != 0 {
		buf = append(buf, ' ')
	}

	a.Value = a.Value.Resolve()
	switch a.Value.Kind() {
	case slog.KindTime:
		buf = fmt.Appendf(buf, "%s=%s", color.CyanString(prefix+a.Key), a.Value.Time().Format("15:04:05 2006-01-02"))
	case slog.KindString:
		buf = fmt.Appendf(buf, "%s=%q", color.CyanString(prefix+a.Key), a.Value)
	case slog.KindGroup:
		attrs := a.Value.Group()
		if len(attrs) == 0 {
			return buf
		}
		if a.Key != "" {
			prefix = fmt.Sprintf("%s%s.", h.preformattedGroups, a.Key)
		}
		for _, ga := range attrs {
			buf = h.appendAttr(buf, ga, prefix)
			buf = append(buf, ' ')
		}
	default:
		buf = fmt.Appendf(buf, "%s=%s", color.CyanString(prefix+a.Key), a.Value)
	}

	return buf
}

func (h *ColoredHandler) appendSource(buf []byte, src slog.Source) []byte {
	relativePath, _ := filepath.Rel(h.projectRoot, src.File)
	fileBuf := fmt.Sprintf("%s:%d", relativePath, src.Line)

	buf = append(buf, faintWhite.Sprint(fileBuf)...)
	return buf
}

func (h *ColoredHandler) appendLevel(buf []byte, level slog.Level) []byte {
	switch level {
	case slog.LevelDebug:
		buf = append(buf, color.MagentaString(level.String())...)
	case slog.LevelInfo:
		buf = append(buf, color.BlueString(level.String())...)
	case slog.LevelWarn:
		buf = append(buf, color.YellowString(level.String())...)
	case slog.LevelError:
		buf = append(buf, color.RedString(level.String())...)
	default:
		buf = append(buf, level.String()...)
	}
	return buf
}

func getProjectRoot() (string, error) {
	path, err := filepath.Abs(".")
	if err != nil {
		return "", fmt.Errorf("getting absolute path: %w", err)
	}

	for path != "/" {
		entries, err := os.ReadDir(path)
		if err != nil {
			return "", fmt.Errorf("reading directory %s: %w", path, err)
		}

		for _, e := range entries {
			if e.Name() == "go.mod" {
				return path, nil
			}
		}
		path = filepath.Dir(path)
	}

	return "", errors.New("didn't find go.mod file in any parent directory")
}
