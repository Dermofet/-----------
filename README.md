http://localhost:8000/docs/index.html - документация API

# Запуск контейнеров
cd ./dev/
docker compose up -d --build

# Запуск проекта
go run ./cmd/music-backend-test/

# Создать swagger.json для документации
swag init -g .\cmd\music-backend-test\main.go --parseInternal
