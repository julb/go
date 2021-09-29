package http

import (
	"context"
	"net/http"

	customContext "github.com/julb/go/pkg/context"
	log "github.com/julb/go/pkg/logging"
	"github.com/julb/go/pkg/tracing"
	"github.com/julb/go/pkg/util/identifier"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type MiddlewareInterceptor func(http.ResponseWriter, *http.Request, http.HandlerFunc)
type MiddlewareHandlerFunc http.HandlerFunc

func (cont MiddlewareHandlerFunc) Intercept(mw MiddlewareInterceptor) MiddlewareHandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		mw(writer, request, http.HandlerFunc(cont))
	}
}

type MiddlewareChain []MiddlewareInterceptor

func (chain MiddlewareChain) Handler(handler http.HandlerFunc) http.Handler {
	curr := MiddlewareHandlerFunc(handler)
	for i := len(chain) - 1; i >= 0; i-- {
		mw := chain[i]
		curr = curr.Intercept(mw)
	}
	return http.HandlerFunc(curr)
}

func NewRequestIdInterceptor() MiddlewareInterceptor {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// Get x-request-id value
		requestId := r.Header.Get(HdrXRequestId)

		// If x-request-id is not set, set value.
		if requestId == "" {
			requestId = identifier.Generate()
			r.Header.Set(HdrXRequestId, requestId)
		}

		// Add it to request context
		ctx := r.Context()

		// Add request ID to the context
		ctx = context.WithValue(ctx, customContext.CtxRequestId, requestId)

		// Add logger to the context
		ctxLogger := ctx.Value(customContext.CtxLogger)
		if ctxLogger == nil {
			loggerWithFieldsSet := log.WithRequestId(requestId)
			ctx = context.WithValue(ctx, customContext.CtxLogger, loggerWithFieldsSet)
		} else {
			loggerWithFieldsSet := ctxLogger.(*log.LogWithContext).WithRequestId(requestId)
			ctx = context.WithValue(ctx, customContext.CtxLogger, loggerWithFieldsSet)
		}

		// Add it to the response.
		w.Header().Add(HdrXRequestId, requestId)

		// Invoke next.
		next(w, r.WithContext(ctx))
	}
}

func NewTrademarkInterceptor() MiddlewareInterceptor {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		// Get x-j3-tm value
		trademark := r.Header.Get(HdrXJ3Tm)

		// If trademark is set
		if trademark != "" {
			ctx := r.Context()

			// Add trademark to the context
			ctx = context.WithValue(ctx, customContext.CtxTrademark, trademark)

			// Add trademark to the logger.
			ctxLogger := ctx.Value(customContext.CtxLogger)
			if ctxLogger == nil {
				loggerWithTrademark := log.WithTrademark(trademark)
				ctx = context.WithValue(ctx, customContext.CtxLogger, loggerWithTrademark)
			} else {
				loggerWithFieldsAndTrademark := ctxLogger.(*log.LogWithContext).WithTrademark(trademark)
				ctx = context.WithValue(ctx, customContext.CtxLogger, loggerWithFieldsAndTrademark)
			}

			next(w, r.WithContext(ctx))
		} else {
			next(w, r)
		}
	}
}

func NewOpentracingInterceptor() MiddlewareInterceptor {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if tracing.IsTracerConfigured() {
			// Create a opentracing context from headers.
			wireContext, err := opentracing.GlobalTracer().Extract(
				opentracing.HTTPHeaders,
				opentracing.HTTPHeadersCarrier(r.Header),
			)
			if err != nil {
				if err == opentracing.ErrSpanContextNotFound {
					log.Trace("no opentracing context received.")
				} else {
					log.Errorf("unable to get opentracing context: %s", err.Error())
				}
			}

			// Get request context
			ctx := r.Context()

			// additional tags
			var requestId string
			var trademark string

			// Add request ID as tag of the span
			ctxRequestId := ctx.Value(customContext.CtxRequestId)
			if ctxRequestId != nil {
				requestId = ctxRequestId.(string)
			}

			// Add trademark to the context
			ctxTrademark := ctx.Value(customContext.CtxTrademark)
			if ctxTrademark != nil {
				trademark = ctxTrademark.(string)
			}

			// start a new span upon this context.
			serverSpan := opentracing.StartSpan(
				r.RequestURI,
				ext.RPCServerOption(wireContext),
				opentracing.Tag{Key: "x-request-id", Value: requestId},
				opentracing.Tag{Key: "tm", Value: trademark},
			)
			defer serverSpan.Finish()

			// Extract trace id and span id
			traceId := tracing.GetTraceId(serverSpan.Context())
			spanId := tracing.GetSpanId(serverSpan.Context())

			// Add logger to the context
			ctxLogger := ctx.Value(customContext.CtxLogger)
			if ctxLogger == nil {
				loggerWithFieldsSet := log.WithTracing(traceId, spanId)
				ctx = context.WithValue(ctx, customContext.CtxLogger, loggerWithFieldsSet)
			} else {
				loggerWithFieldsSet := ctxLogger.(*log.LogWithContext).WithTracing(traceId, spanId)
				ctx = context.WithValue(ctx, customContext.CtxLogger, loggerWithFieldsSet)
			}

			// attach span to request context
			ctx = opentracing.ContextWithSpan(ctx, serverSpan)
			next(w, r.WithContext(ctx))
		} else {
			next(w, r)
		}
	}
}
