# ocp-team-api

## 1. Installation

```
git clone https://github.com/ozoncp/ocp-team-api
```

## 2. Environment setup

### 2.1 Run supported services
```
docker-compose up
```

### 2.2 Apply migrations

```
make migrate
```

## 3. Running

```
make run
```

## 4. Supporting services

### 4.1 Database UI

- 127.0.0.1:8000

### 4.2 Swagger UI

- 127.0.0.1:9080

### 4.3 Metrics

- 127.0.0.1:9100/metrics

### 4.4 Prometheus

- 127.0.0.1:9090

### 4.5 Status server

- 127.0.0.1:9191/health - liveness
- 127.0.0.1:9191/ready - readiness

### 4.6 Kafka UI (through kafdrop)

- 127.0.0.1:9000

### 4.7 Jaeger UI

- 127.0.0.1:16686
