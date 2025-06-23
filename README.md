# Team Collaboration Platform

Многофункциональный backend для платформы командной работы: 
 - управление задачами (Kanban),
 - чат, уведомления,
 - сбор логов, 
 - интеграции через API.

## Архитектура
- Микросервисная структура: каждый сервис отвечает за свою область (users, tasks, chat, notifications, logs и др.).

## Структура проекта
```
GoLessons/
├── .github/                   # Папка для CI/CD 
├── api-gateway/               # Собственный API Gateway сервис
├── migrations/                # SQL-миграции для базы данных (up/down)
├── scripts/                   # Скрипты для инфраструктуры (например, wait-for-it.sh)
├── user-service/              # Микросервис управления пользователями
├── .env                       # Переменные окружения
└── docker-compose.yml         # Docker Compose для запуска всех сервисов и БД
```

## Запуск
```
docker-compose up --build
```

## CI/CD
В проекте реализовано минимальное CI/CD: 
- CI: Реализовано в папке .github
- CD: Реализовано на DockerHub, с помощью секретных ключей на GitHub

## TODO
- OpenAPI/Swagger-документация через gRPC-Gateway
- Helm-чарт для Kubernetes
- Мониторинг (Prometheus-метрики, алерты)
- Централизованный notification-service
- Расширение микросервисов: task-service, chat-service, notification-service, log-service и др.
