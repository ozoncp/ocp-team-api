package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	createSuccessCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "ocp_team_api_create_success_requests",
			Help: "Number of success create requests",
		})
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
)

func Register() {
	prometheus.MustRegister(createSuccessCounter)
	prometheus.MustRegister(updateSuccessCounter)
	prometheus.MustRegister(deleteSuccessCounter)
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
