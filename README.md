# Music test backend
## Описание

Проект реализует API сервиса по хранению музыки пользователей.  

## Требования

Для успешного запуска вашего проекта на Golang необходимо следующее:

- **Golang**: Версия 1.21.0 и выше.

## Установка

1. Клонируйте репозиторий с проектом:
   ```bash
   git clone https://github.com/Dermofet/Music
   ```
   
2. Перейдите в директорию проекта:
    ```bash
    cd Music
    ```
    
3. Установите зависимости:
    ```bash
    go get
    ```
    
## Запуск

Для запуска Docker контейнеров выполните следующие команды:
```bash
./run.sh build_docker
```
Эта команда забилдит проект и запустит его в контейнере.

Для обновления документации Swagger выполните следующую команду:
```bash
swag init -g .\cmd\music-backend-test\main.go --parseInternal
```

Для просмотра документации Swagger перейдите по ссылке:
**http://localhost:8000/docs/index.html**

## Конфигурация

Для конфигурации проекта используется файл **.env**.

#### Переменные окружения

Для настройки http сервера:
- HTTP_HOST=localhost
- HTTP_PORT=8000

Для настройки базы данных:
- DB_HOST=localhost
- DB_PORT=5433
- DB_NAME=devdb
- DB_USER=devuser
- DB_PASS=devpass 
- DB_SSLMODE=disable

## Архитектура базы данных

**СУБД**: PostgreSQL
**Таблицы базы данных:**
- users
  - id (uuid)
  - username (varchar(255))
  - password (varchar(255))
  - role (varchar(255))

- music
  - id (uuid)
  - name (varchar)
  - release_date (date)

- user_music
  - user_id (uuid)
  - music_id (uuid)

## Структура проекта

  - cmd/ 
    - music-backend-test/
      - config/
        - <span style="color: lightblue;">Файлы конфигураций</span>
      - main.go
  - dev/
    - .env
    - docker-compose.yml
  - Dockerfile
  - docs/
    - <span style="color: lightblue;">Файлы документации Swagger</span>
  - go.mod
  - go.sum
  - internal/
    - api/
      - http/
        - handlers/
          - <span style="color: lightblue;">Файлы обработчиков запросов</span>
        - middlewares/
          - <span style="color: lightblue;">Файлы middlewares</span>
        - presenter/
          - <span style="color: lightblue;">Файлы, для преобразования данных</span>
        - router.go
        - server.go
        - view/
          - <span style="color: lightblue;">Файлы представлений сущностей</span>
    - app/
      - app.go
      - migrate.go
      - migrations/
        - <span style="color: lightblue;">Файлы миграций</span>
    - db/
      - <span style="color: lightblue;">Файлы для работы с базой данных</span>
    - entity/
      - <span style="color: lightblue;">Файлы сущностей</span>
    - repository/
      - <span style="color: lightblue;">Файлы репозиториев</span>
    - usecase/
      - <span style="color: lightblue;">Файлы сервисов</span>
  - README.md
  - run.sh
