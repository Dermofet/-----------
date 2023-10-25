http://localhost:8001/api/v0.0.1/docs/index.html - документация API

# Запуск контейнеров
cd ./dev/
docker compose up -d --build

# Запуск проекта
go run ./cmd/go-test-grpc-http/

# Создать swagger.json для документации
swag init -g .\cmd\go-test-grpc-http\main.go --parseInternal
