package loginmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	UserRepoCreateCalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_create_total",
		Help:      "Total number of Create user calls",
		Namespace: "login",
	})
	UserRepoGetByIDCalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_get_by_id_total",
		Help:      "Total number of GetByID user calls",
		Namespace: "login",
	})
	UserRepoGetByEmailCalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_get_by_email_total",
		Help:      "Total number of GetByEmail user calls",
		Namespace: "login",
	})
	UserRepoUpdateCalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_update_total",
		Help:      "Total number of Update user calls",
		Namespace: "login",
	})
	UserRepoDeleteCalls = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_delete_total",
		Help:      "Total number of Delete user calls",
		Namespace: "login",
	})

	UserRepoCreateLastDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "user_repo_create_last_delay_seconds",
		Help:      "Last Create user call delay in seconds",
		Namespace: "login",
	})
	UserRepoGetByIDLastDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "user_repo_get_by_id_last_delay_seconds",
		Help:      "Last GetByID user call delay in seconds",
		Namespace: "login",
	})
	UserRepoGetByEmailLastDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "user_repo_get_by_email_last_delay_seconds",
		Help:      "Last GetByEmail user call delay in seconds",
		Namespace: "login",
	})
	UserRepoUpdateLastDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "user_repo_update_last_delay_seconds",
		Help:      "Last Update user call delay in seconds",
		Namespace: "login",
	})
	UserRepoDeleteLastDelay = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "user_repo_delete_last_delay_seconds",
		Help:      "Last Delete user call delay in seconds",
		Namespace: "login",
	})

	UserRepoCreateSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "user_repo_create_delay_seconds",
		Help:       "Summary of Create user call delays",
		Namespace:  "login",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})
	UserRepoGetByIDSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "user_repo_get_by_id_delay_seconds",
		Help:       "Summary of GetByID user call delays",
		Namespace:  "login",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})
	UserRepoGetByEmailSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "user_repo_get_by_email_delay_seconds",
		Help:       "Summary of GetByEmail user call delays",
		Namespace:  "login",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})
	UserRepoUpdateSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "user_repo_update_delay_seconds",
		Help:       "Summary of Update user call delays",
		Namespace:  "login",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})
	UserRepoDeleteSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "user_repo_delete_delay_seconds",
		Help:       "Summary of Delete user call delays",
		Namespace:  "login",
		Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.01, 0.99: 0.001},
	})

	UserCreateAlreadyExists = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_create_already_exists_count",
		Help:      "Count of attempts to register already existing user",
		Namespace: "login",
	})

	UserGetByIDNotFound = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_get_by_id_not_found_count",
		Help:      "Count of attempts to find by id non-existing user",
		Namespace: "login",
	})

	UserGetByEmailNotFound = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_get_by_email_not_found_count",
		Help:      "Count of attempts to find by email non-existing user",
		Namespace: "login",
	})

	UserUpdateNotFound = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_update_not_found_count",
		Help:      "Count of attempts to update non-existing user. If this metric > 0 => we have critical error in logic.",
		Namespace: "login",
	})
	UserDeleteNotFound = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "user_repo_delete_not_found_count",
		Help:      "Count of attempts to delete non-existing user. If this metric > 0 => we have critical error in logic.",
		Namespace: "login",
	})
)

func Register() {
	prometheus.MustRegister(
		UserRepoCreateCalls, UserRepoGetByIDCalls, UserRepoGetByEmailCalls, UserRepoUpdateCalls, UserRepoDeleteCalls,
		UserRepoCreateLastDelay, UserRepoGetByIDLastDelay, UserRepoGetByEmailLastDelay, UserRepoUpdateLastDelay, UserRepoDeleteLastDelay,
		UserRepoCreateSummary, UserRepoGetByIDSummary, UserRepoGetByEmailSummary, UserRepoUpdateSummary, UserRepoDeleteSummary,
		UserCreateAlreadyExists, UserGetByIDNotFound, UserGetByEmailNotFound, UserUpdateNotFound, UserDeleteNotFound,
	)
}
