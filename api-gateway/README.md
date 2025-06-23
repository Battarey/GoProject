# API Gateway

## Описание
API Gateway — точка входа для всех внешних HTTP-запросов к микросервисам платформы. 
Реализован на Go как отдельный сервис.

- Проксирует запросы `/user/*` на user-service (gRPC-Gateway или REST endpoint).
- Легко расширяется для маршрутизации к другим сервисам (task, chat и др.).
- В перспективе: поддержка JWT, CORS, healthcheck, логирование, rate limiting.

## Запуск
API Gateway будет доступен на http://localhost:8080

## Пример запроса
```
curl http://localhost:8080/user/health
```

## Структура
api-gateway/  
├── main.go                # точка входа, маршрутизация, reverse proxy
└── Dockerfile             # сборка и запуск сервиса

## TODO
- JWT middleware (аутентификация)
- CORS
- Healthcheck endpoint
- Проксирование к другим микросервисам
- Логирование и мониторинг
