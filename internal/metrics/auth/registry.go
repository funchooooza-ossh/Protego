package authmetrics

// Auth usecase metrics
import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// Counters
	LruHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "access_lru_hits_total",
		Help:      "Total number of LRU cache hits",
		Namespace: "auth",
	})
	LruMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "access_lru_misses_total",
		Help:      "Total number of LRU cache misses",
		Namespace: "auth",
	})
	RedisHits = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "access_redis_hits_total",
		Help:      "Total number of Redis cache hits",
		Namespace: "auth",
	})
	RedisMisses = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "access_redis_misses_total",
		Help:      "Total number of Redis cache misses",
		Namespace: "auth",
	})
	DelegateRequests = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "access_database_requests_total",
		Help:      "Total number of requests that reached delegate",
		Namespace: "auth",
	})

	// Delay gauges (last recorded delay in seconds)
	LruLastDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "access_lru_last_delay_seconds",
		Help:      "Last LRU cache access duration in seconds",
		Namespace: "auth",
	})
	RedisLastDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "access_redis_last_delay_seconds",
		Help:      "Last Redis access duration in seconds",
		Namespace: "auth",
	})
	DelegateLastDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "access_delegate_last_delay_seconds",
		Help:      "Last delegate (DB) access duration in seconds",
		Namespace: "auth",
	})

	// Delay summaries (avg, quantiles)
	LruDelaySummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "access_lru_delay_seconds",
		Help:       "LRU cache access time summary",
		Namespace:  "auth",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})
	RedisDelaySummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "access_redis_delay_seconds",
		Help:       "Redis access time summary",
		Namespace:  "auth",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})
	DelegateDelaySummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "access_delegate_delay_seconds",
		Help:       "Delegate (DB) access time summary",
		Namespace:  "auth",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})
)

func Register() {
	prometheus.MustRegister(
		LruHits, LruMisses,
		RedisHits, RedisMisses,
		DelegateRequests,

		LruLastDelay, RedisLastDelay, DelegateLastDelay,
		LruDelaySummary, RedisDelaySummary, DelegateDelaySummary,
	)
}
