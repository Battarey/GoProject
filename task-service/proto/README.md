# proto

Папка содержит gRPC-протоколы и сгенерированные файлы.
Используется для определения API сервиса и генерации кода для взаимодействия между сервисами.

## Структура
proto/
├── task.proto         # описание gRPC API для задач
├── task.pb.go         # сгенерированный Go-код
└── task_grpc.pb.go    # сгенерированный Go-код для gRPC

## Команда для генерации
```
protoc --go_out=paths=source_relative:. --go-grpc_out=paths=source_relative:. proto/task.proto
```