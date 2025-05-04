package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

// Define a metric
var RequestCounter = prometheus.NewCounter(
	prometheus.CounterOpts{
		Name: "tgn_webhook_requests_total",
		Help: "Total number of webhook requests",
	},
)

// New metrics for bot flows
var (
	ErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tgn_errors_total",
			Help: "Total number of errors encountered in the bot",
		},
	)
	ReplyCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tgn_replies_total",
			Help: "Total number of replies sent by the bot",
		},
	)
	StateTransitionCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tgn_state_transitions_total",
			Help: "Total number of state transitions in the bot",
		},
	)
	ConnectAttemptCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tgn_connect_attempts_total",
			Help: "Total number of connection attempts",
		},
	)
	ResetCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tgn_resets_total",
			Help: "Total number of /reset commands received",
		},
	)
	RoleSelectCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tgn_role_select_total",
			Help: "Total number of role selections",
		},
	)
	SalaryParseErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tgn_salary_parse_errors_total",
			Help: "Total number of salary parse errors",
		},
	)
	ResultSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tgn_result_success_total",
			Help: "Total number of successful result calculations",
		},
	)
	ResultErrorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "tgn_result_error_total",
			Help: "Total number of result calculation errors (no overlap)",
		},
	)
)

func init() {
	log.Info("📊 Registering metrics")
	prometheus.MustRegister(RequestCounter)
	prometheus.MustRegister(ErrorCounter)
	prometheus.MustRegister(ReplyCounter)
	prometheus.MustRegister(StateTransitionCounter)
	prometheus.MustRegister(ConnectAttemptCounter)
	prometheus.MustRegister(ResetCounter)
	prometheus.MustRegister(RoleSelectCounter)
	prometheus.MustRegister(SalaryParseErrorCounter)
	prometheus.MustRegister(ResultSuccessCounter)
	prometheus.MustRegister(ResultErrorCounter)
}
