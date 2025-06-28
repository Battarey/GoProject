package handler

import (
	context "context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// HealthCheck реализует gRPC healthcheck endpoint
func (s *TaskServer) HealthCheck(ctx context.Context, _ *struct{}) (*struct{}, error) {
	return &struct{}{}, nil
}

// HealthCheckError возвращает ошибку для проверки мониторинга
func (s *TaskServer) HealthCheckError(ctx context.Context, _ *struct{}) (*struct{}, error) {
	return nil, status.Error(codes.Unavailable, "service unavailable")
}
