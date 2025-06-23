# handler

Папка содержит gRPC-обработчики (handlers) и вспомогательные компоненты для микросервиса user-service.

## Структура папки handler
handler/
├── auth.go           # обработчики регистрации, логина, email, сброса пароля, rate limiting
├── user.go           # CRUD-пользователя (профиль, обновление, удаление, листинг)
├── email.go          # email-логика, генерация токенов, mock email для тестов
├── validation.go     # функции валидации входных данных
├── rate_limiter.go   # in-memory rate limiting
├── utils.go          # вспомогательные функции
└── server.go         # структура UserServer (gRPC-сервер)

Handler-слой организует точки входа (endpoint) gRPC, реализует бизнес-логику, валидацию, защиту и взаимодействие с репозиторием.
