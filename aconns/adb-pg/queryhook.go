package adb_pg

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/fatih/color"
	"github.com/uptrace/bun"
	"io"
	"os"
	"reflect"
	"time"
)

type PGQueryHook struct {
	IsEnabled bool `json:"isEnabled,omitempty"`
	IsVerbose bool `json:"isVerbose,omitempty"`
}

type QueryHookOption func(*QueryHook)

// QueryHookOptionWithEnabled enables/disables the hook.
func QueryHookOptionWithEnabled(on bool) QueryHookOption {
	return func(h *QueryHook) {
		h.enabled = on
	}
}

// QueryHookOptionWithVerbose configures the hook to log all queries
// (by default, only failed queries are logged).
func QueryHookOptionWithVerbose(on bool) QueryHookOption {
	return func(h *QueryHook) {
		h.verbose = on
	}
}

// QueryHookOptionWithWriter sets the log output to an io.Writer
// the default is os.Stderr
func QueryHookOptionWithWriter(w io.Writer) QueryHookOption {
	return func(h *QueryHook) {
		h.writer = w
	}
}

// FromEnv configures the hook using the environment variable value.
// For example, WithEnv("BUNDEBUG"):
//    - BUNDEBUG=0 - disables the hook.
//    - BUNDEBUG=1 - enables the hook.
//    - BUNDEBUG=2 - enables the hook and verbose mode.
//func FromEnv(key string) QueryHookOption {
//	if key == "" {
//		key = "BUNDEBUG"
//	}
//	return func(h *QueryHook) {
//		if env, ok := os.LookupEnv(key); ok {
//			h.enabled = env != "" && env != "0"
//			h.verbose = env == "2"
//		}
//	}
//}

type QueryHook struct {
	enabled bool
	verbose bool
	writer  io.Writer
}

var _ bun.QueryHook = (*QueryHook)(nil)

func NewQueryHook(opts ...QueryHookOption) *QueryHook {
	h := &QueryHook{
		enabled: true,
		writer:  os.Stderr,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (h *QueryHook) BeforeQuery(
	ctx context.Context, event *bun.QueryEvent,
) context.Context {
	return ctx
}

func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	if !h.enabled {
		return
	}

	if !h.verbose {
		switch event.Err {
		case nil, sql.ErrNoRows, sql.ErrTxDone:
			return
		}
	}

	now := time.Now()
	dur := now.Sub(event.StartTime)

	args := []interface{}{
		"[bun]",
		now.Format(" 15:04:05.000 "),
		formatOperation(event),
		fmt.Sprintf(" %10s ", dur.Round(time.Microsecond)),
		event.Query,
	}

	if event.Err != nil {
		typ := reflect.TypeOf(event.Err).String()
		args = append(args,
			"\t",
			color.New(color.BgRed).Sprintf(" %s ", typ+": "+event.Err.Error()),
		)
	}

	fmt.Fprintln(h.writer, args...)
}

func formatOperation(event *bun.QueryEvent) string {
	operation := event.Operation()
	return operationColor(operation).Sprintf(" %-16s ", operation)
}

func operationColor(operation string) *color.Color {
	switch operation {
	case "SELECT":
		return color.New(color.BgGreen, color.FgHiWhite)
	case "INSERT":
		return color.New(color.BgBlue, color.FgHiWhite)
	case "UPDATE":
		return color.New(color.BgYellow, color.FgHiBlack)
	case "DELETE":
		return color.New(color.BgMagenta, color.FgHiWhite)
	default:
		return color.New(color.BgWhite, color.FgHiBlack)
	}
}
