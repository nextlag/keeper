package l

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/fatih/color"

	config "github.com/nextlag/keeper/config/server"
)

type LoggerOptions struct {
	Attrs       []Attr
	SlogOpts    *HandlerOptions
	ProjectPath string
	LogToFile   bool
	LogPath     string
	Line        int
	Out         io.Writer
	log         *log.Logger
}

type logHandler struct {
	opts *LoggerOptions
	Handler
}

var (
	pool = sync.Pool{
		New: func() interface{} {
			return make([]byte, 0, 1024)
		},
	}
)

func NewLogger(cfg *config.Config) *Logger {
	opts := LoggerOptions{
		SlogOpts: &HandlerOptions{
			Level:     cfg.Log.Level,
			AddSource: true,
		},
		ProjectPath: cfg.Log.ProjectPath,
		LogToFile:   cfg.Log.LogToFile,
		LogPath:     cfg.Log.LogPath,
	}
	if opts.LogToFile {
		file, err := os.OpenFile(opts.LogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		opts.Out = file
	} else {
		opts.Out = os.Stdout
	}

	handler := opts.slogHandler(opts.Out, cfg.Log.ProjectPath)
	return slog.New(handler)
}

func (lo *LoggerOptions) slogHandler(out io.Writer, projectPath string) *logHandler {
	h := &logHandler{
		Handler: slog.NewJSONHandler(out, lo.SlogOpts),
		opts: &LoggerOptions{
			ProjectPath: projectPath,
			LogPath:     lo.LogPath,
			LogToFile:   lo.LogToFile,
			log:         log.New(out, "", 0)},
	}

	return h
}

func (lh *logHandler) Handle(_ context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case LevelDebug:
		level = color.YellowString(level)
	case LevelInfo:
		level = color.GreenString(level)
	case LevelWarn:
		level = color.MagentaString(level)
	case LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())

	r.Attrs(func(a Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range lh.opts.Attrs {
		fields[a.Key] = a.Value.Any()
	}

	var buf []byte
	poolBuf := pool.Get()
	switch v := poolBuf.(type) {
	case []byte:
		buf = v
	case *[]byte:
		buf = *v
	default:
		return fmt.Errorf("unexpected type: %T", poolBuf)
	}

	defer func() {
		// Resetting the buffer length to zero before returning to the pool
		buf = buf[:0]
		pool.Put(&buf)
	}()

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.Marshal(fields)
		if err != nil {
			return err
		}
	}

	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "???" // If information is not available
		line = 0
	}
	relPath, err := filepath.Rel(lh.opts.ProjectPath, file)
	if err != nil {
		// If an error occurred, use the full path
		relPath = file
	}

	// Add call information to the log
	lh.opts.LogPath = relPath
	lh.opts.Line = line

	timeStr := r.Time.Format("[02.01.2006 15:04:05.000]")
	msg := color.CyanString(r.Message)

	if lh.opts.LogToFile {
		lh.opts.log.Printf(
			"%s %s %s %s:%d %s\n",
			timeStr,
			r.Level.String()+":",
			r.Message,
			lh.opts.LogPath,
			lh.opts.Line,
			string(b),
		)
		return nil
	}

	lh.opts.log.Println(
		timeStr,
		level,
		msg,
		color.WhiteString("%s:%d", lh.opts.LogPath, lh.opts.Line),
		color.WhiteString(string(b)),
	)

	return nil
}

func (lh *logHandler) WithAttrs(attrs []Attr) Handler {
	return &logHandler{
		Handler: lh.Handler,
		opts: &LoggerOptions{
			Attrs: attrs,
			log:   lh.opts.log,
		},
	}
}

func (lh *logHandler) WithGroup(name string) Handler {
	// TODO: implement
	return &logHandler{
		Handler: lh.Handler.WithGroup(name),
		opts: &LoggerOptions{
			log: lh.opts.log,
		},
	}
}

func L(ctx context.Context) *Logger {
	return loggerFromContext(ctx)
}
