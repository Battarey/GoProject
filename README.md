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
├── scripts/                   # Скрипты для инфраструктуры (например, wait-for-it.sh)
├── e2e_test/                  # Папка для тестов между сервисами
├── task-service/              # Микросервис для задач
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
├── task-service/              # Микросервис для задач
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

## Автоматическая генерация gRPC proto-файлов
Сгенерированные Go-файлы для gRPC автоматически создаются при сборке Docker-образов сервисов. 

## TODO
- OpenAPI/Swagger-документация через gRPC-Gateway
- Helm-чарт для Kubernetes
- Мониторинг (Prometheus-метрики, алерты)
- Централизованный notification-service
- Расширение микросервисов: task-service, chat-service, notification-service, log-service и др.
