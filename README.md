REDME was written by ChatGPT on 2 languages, english and russian, russian version after english

Threat Intelligence API

A RESTful Threat Intelligence API written in Go. The application provides endpoints for creating, retrieving, and updating threat intelligence records using MD5 hashes as unique identifiers.

The project uses a layered architecture and PostgreSQL for persistent data storage.

Tech Stack

* Language: Go
* HTTP: net/http
* Database: PostgreSQL 17
* PostgreSQL Driver: pgx/v5
* Migrations: golang-migrate
* Logging: log/slog
* Containerization: Docker, Docker Compose
* CI: GitHub Actions

Architecture

The application follows a layered architecture:

HTTP Request
     |
     v
  Handler
     |
     v
  Service
     |
     v
Repository Interface
     |
     v
PostgreSQL Repository
     |
     v
 PostgreSQL

Project Structure

threat-intelligence/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── domain/
│   ├── handler/
│   │   └── middleware/
│   ├── repository/
│   │   └── postgres/
│   └── service/
├── api/
│   └── openapi.yaml
├── migrations/
├── tests/
│   └── integration/
├── deployments/
│   └── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md

API Endpoints

All endpoints require an API key passed through the X-API-Key header.

Method	Endpoint	Description
POST	/info	Create a threat record
GET	/info/{md5}	Retrieve a threat by MD5
PUT	/info/{md5}	Update an existing threat

Authentication

Include the API key in every request:

X-API-Key: keyPrivet (this key for test)

Create a Threat

POST /info

Request body:

{
  "md5": "0123456789abcdef0123456789abcdef",
  "name": "Example Trojan",
  "type": "trojan",
  "severity": "high",
  "tags": ["malware", "trojan"]
}

Successful response: 201 Created

{
  "md5": "0123456789abcdef0123456789abcdef",
  "name": "Example Trojan"
}

Get a Threat

GET /info/{md5}

Example:

curl -X GET http://localhost:8080/info/0123456789abcdef0123456789abcdef \
  -H "X-API-Key: your-api-key"

Successful response: 200 OK

Update a Threat

PUT /info/{md5}

Example:

curl -X PUT http://localhost:8080/info/0123456789abcdef0123456789abcdef \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Trojan",
    "type": "trojan",
    "severity": "critical",
    "tags": ["malware", "updated"]
  }'

Successful response: 200 OK

HTTP Status Codes

Status	Description
200	Request completed successfully
201	Threat record created
400	Invalid request or validation error
401	Missing or invalid API key
404	Threat not found
409	Threat with this MD5 already exists
429	Rate limit exceeded
500	Internal server error

Getting Started

Prerequisites

* Docker
* Docker Compose

1. Clone the repository

git clone https://github.com/tPavlovCU/Information_security_REST_Threat.git
cd threat-intelligence

2. Configure environment variables

Create a .env file in the project root:

cp .env.example .env

Configure the values:

HTTP_PORT=8080
DATABASE_URL=postgres://threat_app:APP_DB_PASSWORD@postgres:5432/threats_db?sslmode=disable
HTTP_READ_TIMEOUT=5s
HTTP_WRITE_TIMEOUT=10s
API_KEY=your-secure-api-key
RATE_LIMIT=100
RATE_WINDOW=1m
POSTGRES_PASSWORD=your-postgres-password
APP_DB_PASSWORD=your-app-db-password

Use secure passwords and an API key in your local environment. Do not commit .env to version control.

3. Start the application

docker compose up -d --build

Check running containers:

docker compose ps

View application logs:

docker compose logs -f api

The API will be available at:

http://localhost:8080

4. Stop the application

docker compose down

To stop containers without deleting the PostgreSQL data volume, do not use the -v option.

Database Migrations

Database schema migrations are stored in the migrations/ directory.

Migrations can be applied using golang-migrate:

migrate -path migrations -database "$DATABASE_URL" up

To roll back the latest migration:

migrate -path migrations -database "$DATABASE_URL" down 1

Testing

The project includes unit tests and PostgreSQL integration tests.

Run all tests:

go test ./...

Run tests with verbose output:

go test -v ./...

Run unit tests:

go test ./internal/...

Run integration tests locally:

DATABASE_URL="postgres://threat_user:secret@localhost:5432/threats_test_db?sslmode=disable" \
go test ./tests/integration/...

Integration tests require a running PostgreSQL instance with the project migrations applied.

Important: Integration tests truncate the threats table. Use a dedicated test database, never a production database.

Continuous Integration

GitHub Actions automatically runs checks on pushes and pull requests targeting the main branch.

The CI workflow performs:

* Go formatting checks (gofmt).
* Static analysis (go vet).
* Unit and integration tests.
* Application build.
* PostgreSQL service startup and database migrations.

API Specification

The OpenAPI specification is available at:

api/openapi.yaml

It describes the API endpoints, request schemas, response formats, and authentication requirements.

Security

The application includes:

* API key authentication.
* Constant-time API key comparison.
* IP-based rate limiting.
* HTTP security headers.
* Request ID tracking.
* Structured request logging.
* Panic recovery.
* Restricted PostgreSQL application user.
* Configurable HTTP read and write timeouts.

License

This project is currently available for educational and portfolio purposes.

⸻

Threat Intelligence API

RESTful API для работы с данными об угрозах информационной безопасности, написанный на Go. Приложение позволяет создавать, получать и обновлять записи об угрозах, используя MD5-хеш в качестве уникального идентификатора.

Проект построен на основе многослойной архитектуры и использует PostgreSQL для постоянного хранения данных.

Технологический стек

* Язык: Go
* HTTP: net/http
* База данных: PostgreSQL 17
* Драйвер PostgreSQL: pgx/v5
* Миграции: golang-migrate
* Логирование: log/slog
* Контейнеризация: Docker, Docker Compose
* CI: GitHub Actions

Архитектура

Приложение использует многослойную архитектуру:

HTTP Request
     |
     v
  Handler
     |
     v
  Service
     |
     v
Repository Interface
     |
     v
PostgreSQL Repository
     |
     v
 PostgreSQL

Каждый слой отвечает за отдельную часть приложения:

* Handler — обработка HTTP-запросов, JSON и HTTP-ответов.
* Service — бизнес-логика и валидация данных.
* Repository — интерфейс доступа к данным и реализация работы с PostgreSQL.
* Domain — доменные модели и ошибки.
* Database — создание и управление подключением к базе данных.

Структура проекта

threat-intelligence/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   ├── database/
│   ├── domain/
│   ├── handler/
│   │   └── middleware/
│   ├── repository/
│   │   └── postgres/
│   └── service/
├── api/
│   └── openapi.yaml
├── migrations/
├── tests/
│   └── integration/
├── deployments/
│   └── Dockerfile
├── docker-compose.yml
├── go.mod
└── README.md

API Endpoints

Для доступа ко всем endpoint необходимо передавать API-ключ через заголовок X-API-Key.

Метод	Endpoint	Описание
POST	/info	Создание записи об угрозе
GET	/info/{md5}	Получение информации об угрозе по MD5
PUT	/info/{md5}	Обновление существующей записи

Аутентификация

Добавьте API-ключ в заголовок каждого запроса:

X-API-Key: keyPrivet (это ключ для тестирования)

Создание записи об угрозе

POST /info

Пример тела запроса:

{
  "md5": "0123456789abcdef0123456789abcdef",
  "name": "Example Trojan",
  "type": "trojan",
  "severity": "high",
  "tags": ["malware", "trojan"]
}

Успешный ответ: 201 Created

{
  "md5": "0123456789abcdef0123456789abcdef",
  "name": "Example Trojan"
}

Получение информации об угрозе

GET /info/{md5}

Пример запроса:

curl -X GET http://localhost:8080/info/0123456789abcdef0123456789abcdef \
  -H "X-API-Key: your-api-key"

Успешный ответ: 200 OK.

Обновление записи об угрозе

PUT /info/{md5}

Пример запроса:

curl -X PUT http://localhost:8080/info/0123456789abcdef0123456789abcdef \
  -H "X-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Trojan",
    "type": "trojan",
    "severity": "critical",
    "tags": ["malware", "updated"]
  }'

Успешный ответ: 200 OK.

HTTP-коды ответов

Код	Описание
200	Запрос успешно обработан
201	Запись об угрозе создана
400	Некорректный запрос или ошибка валидации
401	Отсутствует или неверен API-ключ
404	Запись об угрозе не найдена
409	Запись с таким MD5 уже существует
429	Превышен лимит запросов
500	Внутренняя ошибка сервера

Запуск проекта

Необходимые инструменты

* Docker
* Docker Compose

1. Клонирование репозитория

git clone https://github.com/tPavlovCU/Information_security_REST_Threat.git
cd threat-intelligence

2. Настройка переменных окружения

Создайте .env на основе шаблона:

cp .env.example .env

Настройте значения переменных:

HTTP_PORT=8080
DATABASE_URL=postgres://threat_app:APP_DB_PASSWORD@postgres:5432/threats_db?sslmode=disable
HTTP_READ_TIMEOUT=5s
HTTP_WRITE_TIMEOUT=10s
API_KEY=your-secure-api-key
RATE_LIMIT=100
RATE_WINDOW=1m
POSTGRES_PASSWORD=your-postgres-password
APP_DB_PASSWORD=your-app-db-password

Используйте надёжные пароли и API-ключ. Не добавляйте файл .env в Git.

3. Запуск приложения

docker compose up -d --build

Проверить состояние контейнеров:

docker compose ps

Посмотреть логи API:

docker compose logs -f api

После запуска API будет доступен по адресу:

http://localhost:8080

4. Остановка приложения

docker compose down

Для остановки контейнеров с сохранением данных PostgreSQL не используйте параметр -v.

Миграции базы данных

SQL-миграции находятся в директории migrations/.

Для применения миграций используется golang-migrate:

migrate -path migrations -database "$DATABASE_URL" up

Для отмены последней миграции:

migrate -path migrations -database "$DATABASE_URL" down 1

Тестирование

Проект включает unit-тесты и интеграционные тесты с PostgreSQL.

Запуск всех тестов:

go test ./...

Запуск тестов с подробным выводом:

go test -v ./...

Запуск unit-тестов:

go test ./internal/...

Запуск интеграционных тестов локально:

DATABASE_URL="postgres://threat_user:secret@localhost:5432/threats_test_db?sslmode=disable" \
go test ./tests/integration/...

Для интеграционных тестов необходим запущенный PostgreSQL и применённые миграции.

Важно: интеграционные тесты очищают таблицу threats. Используйте отдельную тестовую базу данных, а не production-базу.

Непрерывная интеграция (CI)

GitHub Actions автоматически запускает проверки при push и pull request в ветку main.

Workflow выполняет:

* Проверку форматирования Go-кода (gofmt).
* Статический анализ (go vet).
* Запуск unit- и интеграционных тестов.
* Сборку приложения.
* Запуск PostgreSQL и применение миграций для тестов.

OpenAPI спецификация

OpenAPI-спецификация находится в файле:

api/openapi.yaml

В ней описаны API endpoints, схемы запросов и ответов, а также требования к аутентификации.

Безопасность

В приложении реализованы:

* Аутентификация с помощью API-ключа.
* Constant-time сравнение API-ключей.
* Ограничение количества запросов по IP.
* HTTP security headers.
* Генерация Request ID.
* Структурированное логирование запросов.
* Восстановление после panic.
* Отдельная PostgreSQL-роль приложения с ограниченными правами.
* Настраиваемые таймауты чтения и записи HTTP-сервера.



