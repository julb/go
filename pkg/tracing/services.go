package tracing

import (
	"io"

	"github.com/julb/go/pkg/logging"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-client-go/zipkin"
	"github.com/uber/jaeger-lib/metrics/prometheus"
)

type TracingOpts struct {
	ServiceName string `yaml:"serviceName" json:"serviceName"`
}

var (
	tracingContextCloser io.Closer
)

func Configure(opts *TracingOpts) {
	// Sample configuration for testing. Use constant sampling to sample every trace
	// and enable LogSpan to log every span via configured Logger.
	cfg := jaegerconfig.Configuration{
		ServiceName: opts.ServiceName,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			LogSpans:          true,
			CollectorEndpoint: "http://localhost:14268/api/traces",
		},
	}

	// Zipkin shares span ID between client and server spans; it must be enabled via the following option.
	zipkinPropagator := zipkin.NewZipkinB3HTTPHeaderPropagator()
	metricsFactory := prometheus.New()

	tracingContextCloser, _ = cfg.InitGlobalTracer(
		cfg.ServiceName,
		jaegerconfig.Logger(&jaegerLoggerProxy{
			logger: logging.WithEmptyContext(),
		}),
		jaegerconfig.Metrics(metricsFactory),
		jaegerconfig.Injector(opentracing.HTTPHeaders, zipkinPropagator),
		jaegerconfig.Extractor(opentracing.HTTPHeaders, zipkinPropagator),
		jaegerconfig.ZipkinSharedRPCSpan(false),
	)
}

func IsTracerConfigured() bool {
	return opentracing.IsGlobalTracerRegistered()
}

func GetTraceId(spanContext opentracing.SpanContext) string {
	jaegerSpanContext := spanContext.(jaeger.SpanContext)
	return jaegerSpanContext.TraceID().String()
}

func GetSpanId(spanContext opentracing.SpanContext) string {
	jaegerSpanContext := spanContext.(jaeger.SpanContext)
	return jaegerSpanContext.SpanID().String()
}

func CloseTracingContext() {
	if tracingContextCloser != nil {
		tracingContextCloser.Close()
	}
}

type jaegerLoggerProxy struct {
	logger *logging.LogWithContext
}

func (proxy *jaegerLoggerProxy) Error(msg string) {
	proxy.logger.Error(msg)
}

func (proxy *jaegerLoggerProxy) Infof(msg string, args ...interface{}) {
	proxy.logger.Debugf(msg, args...)
}

func (proxy *jaegerLoggerProxy) Debugf(msg string, args ...interface{}) {
	proxy.logger.Debugf(msg, args...)
}
