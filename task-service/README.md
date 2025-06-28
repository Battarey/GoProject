# task-service

Микросервис управления задачами (Kanban, ToDo, project management).

- Реализован на Go, gRPC, PostgreSQL, JWT, Docker, миграции через golang-migrate.
- Поддержка CRUD задач, смены статуса, фильтрации, тестов, CI/CD, публикации Docker-образа.

## Структура папки
task-service/
├── config                 # Конфигурация сервиса
├── handler                # gRPC-обработчики (endpoint-логика)
├── model                  # Модели данных (структуры задач)
├── proto                  # gRPC-протоколы и сгенерированные файлы
├── repository             # Слой доступа к данным (работа с БД)
├── security               # Логика безопасности (JWT, авторизация)
├── test                   # Модульные тесты для сервиса
├── Dockerfile
├── go.mod
├── go.sum
└── main.go

## API

### gRPC методы
| Метод         | Описание                | Вход/выход                | Ошибки                       |
|---------------|-------------------------|---------------------------|------------------------------|
| CreateTask    | Создать задачу          | CreateTaskRequest/Response| InvalidArgument, Unauth      |
| GetTask       | Получить задачу         | GetTaskRequest/Response   | NotFound, Unauth             |
| UpdateTask    | Обновить задачу         | UpdateTaskRequest/Response| NotFound, PermissionDenied   |
| DeleteTask    | Удалить задачу          | DeleteTaskRequest/Response| NotFound, PermissionDenied   |
| ListTasks     | Список задач            | ListTasksRequest/Response | -                            |
| ChangeStatus  | Сменить статус задачи   | ChangeStatusRequest/Resp  | NotFound, PermissionDenied   |
| HealthCheck   | Проверка статуса        | HealthCheckRequest/Resp   | -                            |

### Пример gRPC-запроса (grpcurl)
```sh
grpcurl -d '{"title":"Test task"}' -H 'authorization: Bearer <JWT>' \
  -plaintext localhost:50052 task.TaskService/CreateTask
```

### Структура Task (protobuf)
```proto
message Task {
  string id = 1;
  string title = 2;
  string description = 3;
  string status = 4; // todo, in_progress, done, ...
  string assignee_id = 5;
  string creator_id = 6;
  string due_date = 7;
  repeated string labels = 8;
  string created_at = 9;
  string updated_at = 10;
}
```

### Авторизация
- Для всех методов (кроме HealthCheck) требуется JWT в metadata:
  - `authorization: Bearer <token>`
- Только создатель, исполнитель или admin может изменять/удалять задачу.

### Ошибки
- `InvalidArgument` — неверные параметры запроса
- `Unauthenticated` — нет или невалидный JWT
- `PermissionDenied` — нет прав на операцию
- `NotFound` — задача не найдена

### Healthcheck
- Метод: `HealthCheck`
- Проверяет доступность сервиса, возвращает статус
