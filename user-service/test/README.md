# Тесты

Папка содержит модульные и интеграционные тесты для микросервиса user-service.

## Структура
test/
├── auth_test.go        # тесты регистрации, логина, email, сброса пароля, rate limiting
├── user_test.go        # тесты CRUD-пользователя
├── email_test.go       # email-моки и edge-cases
├── repository_test.go  # тесты слоя репозитория (работа с БД)
└── testutils.go        # вспомогательные функции для тестов
