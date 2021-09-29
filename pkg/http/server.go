package http

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/julb/go/pkg/crypto/x509"
	log "github.com/julb/go/pkg/logging"
	"github.com/julb/go/pkg/monitoring"
	"github.com/julb/go/pkg/signal"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HttpServerOpts struct {
	Port int `yaml:"port" json:"port"`
	Tls  struct {
		Enabled        bool   `yaml:"enabled" json:"enabled"`
		Certificate    string `yaml:"certificate" json:"certificate"`
		CertificateKey string `yaml:"certificatKey" json:"certificatKey"`
	} `yaml:"tls" json:"tls"`
	Signals []os.Signal
}

type HttpServer struct {
	opts                     *HttpServerOpts
	Router                   *mux.Router
	middlewareChain          MiddlewareChain
	ServerShutdownGracefully <-chan struct{}
}

func DefaultHttpServerOpts() *HttpServerOpts {
	return &HttpServerOpts{
		Port:    8080,
		Signals: []os.Signal{os.Interrupt, syscall.SIGINT, syscall.SIGTERM},
	}
}

func NewHttpServer(opts *HttpServerOpts) (*HttpServer, error) {
	// Add default middlewares
	middlewareChain := MiddlewareChain{
		NewRequestIdInterceptor(),
		NewTrademarkInterceptor(),
		NewOpentracingInterceptor(),
	}

	// Configure runtime information
	registerRuntimeInformation(opts)

	// Register default metrics.
	prometheus.MustRegister(monitoring.NewPrometheusCollector())
	prometheus.MustRegister(NewPrometheusCollector(opts))

	// Add default handlers
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler()).Methods(http.MethodGet)
	router.Handle("/health", middlewareChain.Handler(HealthCheckHandler)).Methods(http.MethodGet)
	router.Handle("/healthz", middlewareChain.Handler(HealthCheckHandler)).Methods(http.MethodGet)
	router.Handle("/info", middlewareChain.Handler(InfoHandler)).Methods(http.MethodGet)
	router.Handle("/runtime", middlewareChain.Handler(RuntimeHandler)).Methods(http.MethodGet)
	http.Handle("/", router)

	// return HTTP server
	return &HttpServer{
		opts:            opts,
		Router:          router,
		middlewareChain: middlewareChain,
	}, nil
}

// Start the server
func (httpServer *HttpServer) Start() {
	// Instantiate server
	var golangHttpServer = &http.Server{
		Addr: fmt.Sprintf(":%d", httpServer.opts.Port),
	}

	// Trap signal
	trapInterruptionSignal := signal.TrapSignal(func(signal os.Signal) {
		// Shutdown the server gracefully
		log.Debug("initiating a graceful shutdown of the server.")
		if err := golangHttpServer.Shutdown(context.Background()); err != nil {
			log.Errorf("error when shutting down server: %v", err)
		}
		log.Info("http server shut down gracefully.")
	}, httpServer.opts.Signals...)

	// Assign context to server object.
	httpServer.ServerShutdownGracefully = trapInterruptionSignal.SignalHandled

	// Starting HTTP server.
	if httpServer.opts.Tls.Enabled {
		log.Infof("start https listener on %s", golangHttpServer.Addr)
		if err := golangHttpServer.ListenAndServeTLS(httpServer.opts.Tls.Certificate, httpServer.opts.Tls.CertificateKey); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	} else {
		log.Infof("start http listener on %s", golangHttpServer.Addr)
		if err := golangHttpServer.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}
}

// register the runtime information in the contributor.
func registerRuntimeInformation(opts *HttpServerOpts) {
	// get tls metadata if any.
	var tlsCertificateMetadata *x509.TlsCertificateMetadata
	if opts.Tls.Enabled {
		tlsCertificateMetadata, _ = x509.ParsePemFileAndGetMetadata(opts.Tls.Certificate)
	}

	// register in contributor.
	monitoring.GetRuntimeContributor().With("server", map[string]interface{}{
		"port": opts.Port,
		"tls": map[string]interface{}{
			"enabled":     opts.Tls.Enabled,
			"certificate": tlsCertificateMetadata,
		},
	})
}
