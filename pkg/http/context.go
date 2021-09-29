package http

import (
	"net/http"

	cc "github.com/julb/go/pkg/context"
	log "github.com/julb/go/pkg/logging"
)

// Returns the contextual logger or a classical one if not set.
func GetCtxLogger(r *http.Request) *log.LogWithContext {
	return cc.GetCtxLogger(r.Context())
}

// Returns the contextual request ID or "" if not set.
func GetCtxRequestId(r *http.Request) string {
	return cc.GetCtxRequestId(r.Context())
}

// Returns the contextual trademark or "" if not set.
func GetCtxTrademark(r *http.Request) string {
	return cc.GetCtxTrademark(r.Context())
}
