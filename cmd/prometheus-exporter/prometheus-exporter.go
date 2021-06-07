package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"log/syslog"

	"github.com/julb/go/pkg/monitoring"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	logger "github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func recordMetrics() {
	go func() {
		for {
			logger.Info("Refresh metrics")
			opsProcessed.Inc()
			time.Sleep(2 * time.Second)
		}
	}()
}

func startConfig() {
	// Env management
	viper.AutomaticEnv()
	viper.SetEnvPrefix("julb")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))

	// Flag management
	pflag.String("messageflag", "", "some custom message from flag")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	// Config file
	viper.SetConfigName("prometheus-exporter") // name of config file (without extension)
	viper.AddConfigPath("/etc/julb")           // path to look for the config file in
	viper.AddConfigPath("$HOME/.julb")         // call multiple times to add many search paths
	viper.AddConfigPath("./configs")           // optionally look for config in the working directory
	viper.AddConfigPath(".")                   // optionally look for config in the working directory
	err := viper.ReadInConfig()                // Find and read the config file
	if err != nil {                            // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
}

func buildLogger() *logger.Entry {
	// Log as JSON instead of the default ASCII formatter.
	logger.SetFormatter(&logger.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logger.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logger.SetLevel(logger.InfoLevel)

	// Add syslog hook
	hook, err := logrus_syslog.NewSyslogHook("udp", "localhost:10514", syslog.LOG_INFO, "golang")
	if err != nil {
		logger.Error("Unable to connect to local syslog daemon")
	} else {
		logger.AddHook(hook)
		logger.Info("Hook added!")
	}

	return logger.WithFields(logger.Fields{"request_id": "aaa", "user_ip": "0.0.0.0"})
}

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
	})
)

func main() {
	startConfig()
	localLogger := buildLogger()

	// Debug
	fmt.Printf("Message from config file: %s\n", viper.GetString("message"))
	fmt.Printf("Hierarchical message from config file: %s\n", viper.GetString("greetings.julien-b"))
	fmt.Printf("Message from env file: %s\n", viper.GetString("messagenv"))
	fmt.Printf("Message from flag: %s\n", viper.GetString("messageflag"))
	fmt.Printf("Dict: %s\n", viper.GetStringMap("info"))

	json.NewEncoder(os.Stdout).Encode(viper.GetStringMap("info"))

	recordMetrics()

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", monitoring.GetHealth)
	http.HandleFunc("/info", monitoring.GetHealth)

	localLogger.Info("Starting server on :2112 port.")

	localLogger.Fatal(http.ListenAndServe(":2112", nil))
}
