package internal

//go:generate mockgen -destination=./mocks/repo_mock.go -package=mocks github.com/ozoncp/ocp-team-api/internal/repo IRepo
//go:generate mockgen -destination=./mocks/flusher_mock.go -package=mocks github.com/ozoncp/ocp-team-api/internal/flusher IFlusher
//go:generate mockgen -destination=./mocks/producer_mock.go -package=mocks github.com/ozoncp/ocp-team-api/internal/kafka IProducer
