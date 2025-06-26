package lifespanmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	ComponentErrors = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "observability",
			Name:      "component_errors_total",
			Help:      "Total number of component errors, labeled by component and error type",
		},
		[]string{"component", "type"}, // e.g., component=redis, type=timeout
	)
)

func Register() {
	prometheus.MustRegister(ComponentErrors)
}

func Inc(component, errType string) {
	ComponentErrors.WithLabelValues(component, errType).Inc()
}
