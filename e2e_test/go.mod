module e2e_test

go 1.23

require (
	google.golang.org/grpc v1.62.1
	task-service/proto v0.0.0
	user-service/proto v0.0.0
)

require (
	github.com/golang/protobuf v1.5.3 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240123012728-ef4313101c80 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
)

replace task-service/proto => ../task-service/proto

replace user-service/proto => ../user-service/proto
