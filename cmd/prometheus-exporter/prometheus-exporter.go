package main

import (
	"github.com/mitchellh/mapstructure"

	"github.com/julb/go/pkg/http"
	log "github.com/julb/go/pkg/logging"
	"github.com/julb/go/pkg/monitoring"
	"github.com/julb/go/pkg/settings"
	"github.com/julb/go/pkg/tracing"
)

// Configure logging
func configureLogging(settings settings.Settings) {
	loggingOpts := log.DefaultLoggingOpts()
	err := mapstructure.Decode(settings["logging"], &loggingOpts)
	if err != nil { // Handle errors reading the config file
		log.Fatal("Unable to decode logging configuration.")
	}
	err = log.Configure(loggingOpts)
	if err != nil {
		log.Fatal("Unable to configure logging.")
	}
}

// Configure HTTP server
func configureHttpServer(settings settings.Settings) *http.HttpServer {
	httpServerOpts := http.DefaultHttpServerOpts()
	err := mapstructure.Decode(settings["server"], &httpServerOpts)
	if err != nil {
		log.Fatal("Unable to decode http server configuration.")
	}
	httpServer, err := http.NewHttpServer(httpServerOpts)
	if err != nil {
		log.Fatal("Unable to create http server configuration.")
	}
	return httpServer
}

func configureTracing(settings settings.Settings) {
	tracingOpts := &tracing.TracingOpts{
		ServiceName: "test",
	}
	tracing.Configure(tracingOpts)
}

func main() {
	settingsVar := settings.ParseAndGet()

	// Configure contributors
	monitoring.GetInfoContributor().WithMap(settings.GetKey("info"))

	// Configure logging.
	configureLogging(settingsVar)

	// Configure tracer
	configureTracing(settingsVar)

	// Configure HTTP server
	var httpServer = configureHttpServer(settingsVar)

	// Start server
	httpServer.Start()

	// Exit program when server has gracefully shutdown.
	<-httpServer.ServerShutdownGracefully

	// Close tracing context
	tracing.CloseTracingContext()
}
