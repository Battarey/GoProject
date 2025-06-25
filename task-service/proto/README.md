# proto

Папка содержит gRPC-протоколы (task.proto) и сгенерированные файлы (task.pb.go, task_grpc.pb.go).

- task.proto: описание сервисов и сообщений для задач (CRUD, смена статуса, фильтрация)

## Команда для инициализации
```
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. proto/task.proto
```