package bundebug

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"utilware/bun"
	"utilware/logger"
)

type Option func(*QueryHook)

// WithEnabled enables/disables the hook.
func WithEnabled(on bool) Option {
	return func(h *QueryHook) {
		h.enabled = on
	}
}

// WithVerbose configures the hook to log all queries
// (by default, only failed queries are logged).
func WithVerbose(on bool) Option {
	return func(h *QueryHook) {
		h.verbose = on
	}
}

// EmptyLine adds an empty line before each query.
func WithEmptyLine(on bool) Option {
	return func(h *QueryHook) {
		h.emptyLine = on
	}
}

// FromEnv configures the hook using the environment variable value.
// For example, WithEnv("BUNDEBUG"):
//   - BUNDEBUG=0 - disables the hook.
//   - BUNDEBUG=1 - enables the hook.
//   - BUNDEBUG=2 - enables the hook and verbose mode.
func FromEnv(keys ...string) Option {
	if len(keys) == 0 {
		keys = []string{"BUNDEBUG"}
	}
	return func(h *QueryHook) {
		for _, key := range keys {
			if env, ok := os.LookupEnv(key); ok {
				h.enabled = env != "" && env != "0"
				h.verbose = env == "2"
				break
			}
		}
	}
}

type QueryHook struct {
	enabled   bool
	verbose   bool
	emptyLine bool
}

var _ bun.QueryHook = (*QueryHook)(nil)

func NewQueryHook(opts ...Option) *QueryHook {
	h := &QueryHook{
		enabled: true,
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

func (h *QueryHook) AfterQuery(ctx context.Context, evt *bun.QueryEvent) {
	if !h.enabled {
		return
	}

	if !h.verbose {
		switch evt.Err {
		case nil, sql.ErrNoRows, sql.ErrTxDone:
			return
		}
	}

	if h.emptyLine {
		defer fmt.Println()
	}

	if evt.Err != nil {
		logger.Printf("[bun] [%s] %s\r\n[%s] %s",
			formatOperation(evt), time.Since(evt.StartTime), evt.Query, evt.Err)
		return
	}

	if evt.Result != nil {
		if evt.Operation() == "SELECT" {
			lastInsertId, _ := evt.Result.LastInsertId()

			logger.Printf("[bun] [%s] %s (%d)\r\n%s",
				formatOperation(evt), time.Since(evt.StartTime), lastInsertId, evt.Query)
			return
		}

		rowsAffected, _ := evt.Result.RowsAffected()

		logger.Printf("[bun] [%s] %s (%d)\r\n %s",
			formatOperation(evt), time.Since(evt.StartTime), rowsAffected, evt.Query)
	}

	logger.Printf("[bun] [%s] %s\r\n%s",
		formatOperation(evt), time.Since(evt.StartTime), evt.Query)
}

func formatOperation(event *bun.QueryEvent) string {
	operation := event.Operation()
	return logger.ANSICode(operationColor(operation), operation)
}

func operationColor(operation string) string {
	switch operation {
	case "SELECT":
		return logger.ANSIgreen
	case "INSERT":
		return logger.ANSIblue
	case "UPDATE":
		return logger.ANSIyellow
	case "DELETE":
		return logger.ANSIred
	default:
		return logger.ANSIwhite
	}
}
