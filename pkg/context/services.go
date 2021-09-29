package context

import (
	"context"

	log "github.com/julb/go/pkg/logging"
)

// Returns the contextual request ID or "" if not set.
func GetCtxRequestId(context context.Context) string {
	value := context.Value(CtxRequestId)
	if value == nil {
		return ""
	}
	return value.(string)
}

// Returns the contextual trademark or "" if not set.
func GetCtxTrademark(context context.Context) string {
	value := context.Value(CtxTrademark)
	if value == nil {
		return ""
	}
	return value.(string)
}

// Returns the contextual logger or a classical one if not set.
func GetCtxLogger(context context.Context) *log.LogWithContext {
	ctxLogger := context.Value(CtxLogger)
	if ctxLogger == nil {
		ctxLogger = log.WithEmptyContext()
	}
	return ctxLogger.(*log.LogWithContext)
}
