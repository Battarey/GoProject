# API Gateway

API Gateway — точка входа для HTTP-запросов к микросервисам Team Collaboration Platform. Реализован на Go, легко расширяется для новых сервисов.

- Reverse proxy для маршрута `/user/*` на user-service
- JWT middleware (аутентификация)
- Rate limiting (ограничение частоты запросов)
- CORS middleware (разрешение кросс-доменных запросов)
- Healthcheck endpoint `/health`
- Готов к расширению (task, chat и др.)

## Структура
```
api-gateway/
├── main.go                # Точка входа, маршрутизация, запуск сервера
├── handlers/              # Обработчики (health, proxy)
│   ├── health.go
│   └── proxy.go
├── middlewares/           # Middleware: JWT, CORS, rate limiting
│   ├── cors.go
│   ├── jwt.go
│   └── ratelimit.go
├── test/                  # Unit-тесты middleware и обработчиков
│   ├── cors_test.go
│   ├── health_test.go
│   ├── jwt_test.go
│   └── ratelimit_test.go
├── Dockerfile             # Сборка и запуск сервиса
└── go.mod                 # Go modules
```

## Примеры запросов
Проверка работоспособности:
```
curl http://localhost:8080/health
```
Запрос к user-service через gateway:
```
curl http://localhost:8080/user/profile
```
