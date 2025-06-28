# test

Папка содержит модульные и интеграционные тесты для микросервиса task-service.

## Структура
```
test/
├── task_create_test.go   # тесты создания задач и валидации
├── task_update_test.go   # тесты обновления задач и edge-cases
├── task_delete_test.go   # тесты удаления задач и edge-cases
├── task_status_test.go   # тесты смены статуса задач
├── task_get_test.go      # тесты получения задач
├── testutils.go          # вспомогательные функции для тестов (setup, JWT, context)
└── README.md             # описание тестов и подходов
```
