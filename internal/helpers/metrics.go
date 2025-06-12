package helpers

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func ObserveDuration(g prometheus.Gauge, s prometheus.Summary, start time.Time) {
	duration := time.Since(start).Seconds()
	g.Set(duration)
	s.Observe(duration)
}
