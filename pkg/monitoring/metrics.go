package monitoring

import (
	"github.com/julb/go/pkg/build"
	"github.com/prometheus/client_golang/prometheus"
)

type PrometheusCollector struct {
	buildInfo *prometheus.Desc
	up        *prometheus.Desc
}

func NewPrometheusCollector() *PrometheusCollector {
	// Get build_info labels.
	buildInfoMetricLabels := prometheus.Labels{
		"group":            build.Info.Group,
		"artifact":         build.Info.Artifact,
		"name":             build.Info.Name,
		"version":          build.Info.Version,
		"time":             build.Info.Time,
		"buildVersion":     build.Info.BuildVersion,
		"gitRevision":      build.Info.GitRevision,
		"gitShortRevision": build.Info.GitShortRevision,
		"arch":             build.Info.Arch,
	}

	// Return prometheus collector.
	return &PrometheusCollector{
		buildInfo: prometheus.NewDesc("build_info", "Current build information", nil, buildInfoMetricLabels),
		up:        prometheus.NewDesc("up", "Health status of the system", nil, nil),
	}
}

// Describe the prometheus collector with all metrics.
func (collector *PrometheusCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.up
	ch <- collector.buildInfo
}

// Collect all metric values.
func (collector *PrometheusCollector) Collect(ch chan<- prometheus.Metric) {
	// Build info.
	ch <- prometheus.MustNewConstMetric(collector.buildInfo, prometheus.GaugeValue, 1)

	// Up
	upValue := mapSystemStatusToMetricValue(GetSystemStatus())
	ch <- prometheus.MustNewConstMetric(collector.up, prometheus.GaugeValue, upValue)
}

// Map system status to a prometheus metric value.
func mapSystemStatusToMetricValue(systemStatus SystemStatus) float64 {
	switch systemStatus {
	case Up:
		return 1
	case Down:
		return 0
	case OutOfService:
		return -1
	case Partial:
		return -2
	case Unknown:
		return -3
	default:
		return -3
	}
}
