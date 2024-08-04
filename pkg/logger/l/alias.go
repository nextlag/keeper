package l

import (
	"fmt"
	"log/slog"
	"runtime"
	"strings"
	"time"
)

const (
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelDebug = slog.LevelDebug
)

type (
	Logger         = slog.Logger
	Attr           = slog.Attr
	Level          = slog.Level
	Handler        = slog.Handler
	Value          = slog.Value
	HandlerOptions = slog.HandlerOptions
	LogValuer      = slog.LogValuer
)

var (
	NewTextHandler = slog.NewTextHandler
	NewJSONHandler = slog.NewJSONHandler
	New            = slog.New
	SetDefault     = slog.SetDefault

	StringAttr   = slog.String
	BoolAttr     = slog.Bool
	Float64Attr  = slog.Float64
	AnyAttr      = slog.Any
	DurationAttr = slog.Duration
	IntAttr      = slog.Int
	Int64Attr    = slog.Int64
	Uint64Attr   = slog.Uint64

	GroupValue = slog.GroupValue
	Group      = slog.Group
)

func Float32Attr(key string, val float32) Attr {
	return slog.Float64(key, float64(val))
}

func UInt32Attr(key string, val uint32) Attr {
	return slog.Int(key, int(val))
}

func Int32Attr(key string, val int32) Attr {
	return slog.Int(key, int(val))
}

func TimeAttr(key string, time time.Time) Attr {
	return slog.String(key, time.String())
}

func ErrAttr(err error) (attr Attr) {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return attr
	}
	fullFuncName := runtime.FuncForPC(pc).Name()
	parts := strings.Split(fullFuncName, "/")
	funcName := parts[len(parts)-1]
	packageAndFunc := strings.Split(funcName, ".")
	packageName := strings.Join(packageAndFunc[:len(packageAndFunc)-1], ".")
	funcName = packageAndFunc[len(packageAndFunc)-1]
	return slog.String(fmt.Sprintf("%s.%s", packageName, funcName), err.Error())
}
