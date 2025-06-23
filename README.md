# Team Collaboration Platform

## Краткое описание
Многофункциональный backend для платформы командной работы: 
 - управление задачами (Kanban),
 - чат, уведомления,
 - сбор логов, 
 - интеграции через API.

## Архитектура
- Микросервисная структура: каждый сервис отвечает за свою область (users, tasks, chat, notifications, logs и др.).
- Взаимодействие между сервисами через gRPC и/или message broker.

## Структура проекта
```
GoLessons/
├── migrations                 # SQL-миграции для базы данных (up/down)
├── scripts                    # Скрипты для инфраструктуры (например, wait-for-it.sh)
├── user-service/              # Основной микросервис управления пользователями
│   ├── config                 # Конфигурация сервиса
│   ├── handler                # gRPC-обработчики (endpoint-логика)
│   ├── model                  # Модели данных (структуры пользователей)
│   ├── proto                  # gRPC-протоколы и сгенерированные файлы
│   ├── repository             # Слой доступа к данным (работа с БД)
│   ├── security               # Логика безопасности (JWT, авторизация)
│   ├── test                   # Модульные тесты для сервиса
│   ├── Dockerfile             # Dockerfile для сборки user-service
│   ├── go.mod
│   ├── go.sum
│   └── main.go
├── .github/                   # Папка для CI/CD 
├── docker-compose.yml         # Docker Compose для запуска всех сервисов и БД
├── .env                       # Переменные окружения
└── README.md                  # Глобальное описание и архитектура проекта
```

## user-service
- Сервис управления пользователями, аутентификацией и ролями.
- Реализован на Go, gRPC, PostgreSQL, JWT, Docker, миграции через golang-migrate.
- Поддержка CRUD, ролей, валидации, тестов, CI/CD, публикации Docker-образа.

## Запуск
```
docker-compose up --build
```

## TODO (сделать позже)
- Нагрузочные тесты (k6, vegeta, autocannon)
- Двухфакторная аутентификация (2FA)
- Аудит действий пользователя (логирование входов, изменений профиля)
- Интеграция с внешними сервисами (email, SMS, OAuth)
- OpenAPI/Swagger-документация через gRPC-Gateway
- Healthcheck endpoint для Kubernetes/DevOps
- Helm-чарт для Kubernetes
- Мониторинг (Prometheus-метрики, алерты)
- Централизованный notification-service
- API Gateway
- Расширение микросервисов: task-service, chat-service, notification-service, log-service и др.
