package http

import (
	"github.com/julb/go/pkg/crypto/x509"
	"github.com/julb/go/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusCollector struct {
	httpServerOpts               *HttpServerOpts
	tlsServerCertificateMetadata *x509.TlsCertificateMetadata

	// Metrics
	tlsServerCertificateRemainingDays *prometheus.Desc
	tlsServerCertificateExpired       *prometheus.Desc
	tlsServerCertificateValid         *prometheus.Desc
}

// Create a prometheus collector for HTTP server
func NewPrometheusCollector(httpServerOpts *HttpServerOpts) *PrometheusCollector {
	collector := &PrometheusCollector{
		httpServerOpts: httpServerOpts,
	}

	if httpServerOpts.Tls.Enabled {
		// parse certificate
		tlsServerCertificateMetadata, err := x509.ParsePemFileAndGetMetadata(httpServerOpts.Tls.Certificate)
		if err != nil {
			logging.Errorf("fail to extract tls server certificate information: %v", err)
		}

		collector.tlsServerCertificateMetadata = tlsServerCertificateMetadata
		collector.tlsServerCertificateRemainingDays = prometheus.NewDesc("tls_server_certificate_expiry_in_days", "Flag indicating the number of remaining days before the TLS server certificate expire", nil, nil)
		collector.tlsServerCertificateExpired = prometheus.NewDesc("tls_server_certificate_expired", "Flag indicating if TLS server certificate is expired", nil, nil)
		collector.tlsServerCertificateValid = prometheus.NewDesc("tls_server_certificate_valid", "Flag indicating if TLS server certificate is valid", nil, nil)
	}

	return collector
}

// Describe the prometheus collector with all metrics.
func (collector *PrometheusCollector) Describe(ch chan<- *prometheus.Desc) {
	if collector.httpServerOpts.Tls.Enabled {
		ch <- collector.tlsServerCertificateRemainingDays
		ch <- collector.tlsServerCertificateExpired
		ch <- collector.tlsServerCertificateValid
	}
}

// Collect all metric values.
func (collector *PrometheusCollector) Collect(ch chan<- prometheus.Metric) {
	if collector.httpServerOpts.Tls.Enabled {
		// Initialize with default value.
		tlsServerCertificateRemainingDaysMetricValue := float64(-1)
		tlsServerCertificateExpiredMetricValue := float64(-1)
		tlsServerCertificateValidMetricValue := float64(-1)

		// parse certificate
		if collector.tlsServerCertificateMetadata != nil {
			tlsServerCertificateRemainingDaysMetricValue = float64(collector.tlsServerCertificateMetadata.Validity.RemainingDays)
			tlsServerCertificateExpiredMetricValue = bfloat64(collector.tlsServerCertificateMetadata.Validity.Expired)
			tlsServerCertificateValidMetricValue = bfloat64(collector.tlsServerCertificateMetadata.Validity.Valid)
		}

		// provision metrics
		ch <- prometheus.MustNewConstMetric(collector.tlsServerCertificateRemainingDays, prometheus.GaugeValue, tlsServerCertificateRemainingDaysMetricValue)
		ch <- prometheus.MustNewConstMetric(collector.tlsServerCertificateExpired, prometheus.GaugeValue, tlsServerCertificateExpiredMetricValue)
		ch <- prometheus.MustNewConstMetric(collector.tlsServerCertificateValid, prometheus.GaugeValue, tlsServerCertificateValidMetricValue)
	}
}

func bfloat64(booleanValue bool) float64 {
	floatValue := float64(0)
	if booleanValue {
		floatValue = 1
	}
	return floatValue
}
