package monitoring

import (
	"github.com/julb/go/pkg/build"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func init() {
	promauto.NewGauge(prometheus.GaugeOpts{
		Name: "build_info",
		ConstLabels: prometheus.Labels{
			"group":            build.Info.Group,
			"artifact":         build.Info.Artifact,
			"name":             build.Info.Name,
			"version":          build.Info.Version,
			"time":             build.Info.Time,
			"buildVersion":     build.Info.BuildVersion,
			"gitRevision":      build.Info.GitRevision,
			"gitShortRevision": build.Info.GitShortRevision,
			"arch":             build.Info.Arch,
		},
	}).Set(1)
}
