package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	createSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ocp_team_api_create_success_requests",
			Help: "Number of success create requests",
		},
	)
	updateSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ocp_team_api_update_success_requests",
			Help: "Number of success update requests",
		},
	)
	deleteSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ocp_team_api_delete_success_requests",
			Help: "Number of success delete requests",
		},
	)
	invalidRequestsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ocp_team_api_invalid_requests",
			Help: "Number of incoming invalid requests",
		},
	)
	totalRequestsCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ocp_team_api_total_requests",
			Help: "Number of total incoming requests",
		},
	)
)

func Register() {
	prometheus.MustRegister(createSuccessCounter)
	prometheus.MustRegister(updateSuccessCounter)
	prometheus.MustRegister(deleteSuccessCounter)

	prometheus.MustRegister(invalidRequestsCounter)
	prometheus.MustRegister(totalRequestsCounter)
}

func IncCreateSuccessCounter() {
	createSuccessCounter.Inc()
}

func IncUpdateSuccessCounter() {
	updateSuccessCounter.Inc()
}

func IncDeleteSuccessCounter() {
	deleteSuccessCounter.Inc()
}

func IncInvalidRequestsCounter() {
	invalidRequestsCounter.Inc()
}

func IncTotalRequestsCounter() {
	totalRequestsCounter.Inc()
}
