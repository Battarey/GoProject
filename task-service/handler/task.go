package handler

import (
	"context"
	"errors"
	"task-service/model"
	pb "task-service/proto"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
)

func (s *TaskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.CreateTaskResponse, error) {
	if err := ValidateCreateTaskInput(req.Title); err != nil {
		return nil, GRPCError(err.Error(), codes.InvalidArgument)
	}
	task := &model.Task{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Status:      "todo",
		AssigneeID:  uuid.Nil,
		CreatorID:   uuid.Nil, // Можно получить из JWT
		Labels:      req.Labels,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if req.AssigneeId != "" {
		if id, err := uuid.Parse(req.AssigneeId); err == nil {
			task.AssigneeID = id
		}
	}
	if req.DueDate != "" {
		if due, err := time.Parse(time.RFC3339, req.DueDate); err == nil {
			task.DueDate = &due
		}
	}
	if err := s.Repo.CreateTask(task); err != nil {
		return nil, err
	}
	return &pb.CreateTaskResponse{TaskId: task.ID.String()}, nil
}

func (s *TaskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.GetTaskResponse, error) {
	task, err := s.Repo.GetTaskByID(req.TaskId)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, errors.New("task not found")
	}
	return &pb.GetTaskResponse{Task: toProtoTask(task)}, nil
}

func (s *TaskServer) UpdateTask(ctx context.Context, req *pb.UpdateTaskRequest) (*pb.UpdateTaskResponse, error) {
	if err := ValidateUpdateTaskInput(req.Title); err != nil {
		return nil, GRPCError(err.Error(), codes.InvalidArgument)
	}
	userID, role, err := GetAuthContext(ctx, s.JwtService)
	if err != nil {
		return nil, GRPCError("unauthorized", codes.Unauthenticated)
	}
	task, err := s.Repo.GetTaskByID(req.TaskId)
	if err != nil || task == nil {
		return nil, GRPCError("task not found", codes.NotFound)
	}
	// Только исполнитель, создатель или админ может обновлять задачу
	if userID != task.AssigneeID.String() && userID != task.CreatorID.String() && role != "admin" {
		return nil, GRPCError("forbidden", codes.PermissionDenied)
	}
	task.Title = req.Title
	task.Description = req.Description
	if req.AssigneeId != "" {
		if id, err := uuid.Parse(req.AssigneeId); err == nil {
			task.AssigneeID = id
		}
	}
	if req.DueDate != "" {
		if due, err := time.Parse(time.RFC3339, req.DueDate); err == nil {
			task.DueDate = &due
		}
	}
	task.Labels = req.Labels
	task.UpdatedAt = time.Now()
	if err := s.Repo.UpdateTask(task); err != nil {
		return nil, err
	}
	return &pb.UpdateTaskResponse{TaskId: task.ID.String()}, nil
}

func (s *TaskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteTaskResponse, error) {
	userID, role, err := GetAuthContext(ctx, s.JwtService)
	if err != nil {
		return &pb.DeleteTaskResponse{Success: false}, GRPCError("unauthorized", codes.Unauthenticated)
	}
	task, err := s.Repo.GetTaskByID(req.TaskId)
	if err != nil || task == nil {
		return &pb.DeleteTaskResponse{Success: false}, GRPCError("task not found", codes.NotFound)
	}
	if userID != task.CreatorID.String() && role != "admin" {
		return &pb.DeleteTaskResponse{Success: false}, GRPCError("forbidden", codes.PermissionDenied)
	}
	if err := s.Repo.DeleteTask(req.TaskId); err != nil {
		return &pb.DeleteTaskResponse{Success: false}, GRPCError("internal error", codes.Internal)
	}
	return &pb.DeleteTaskResponse{Success: true}, nil
}

func (s *TaskServer) ListTasks(ctx context.Context, req *pb.ListTasksRequest) (*pb.ListTasksResponse, error) {
	offset := int(req.Page-1) * int(req.PageSize)
	limit := int(req.PageSize)
	tasks, err := s.Repo.ListTasks(req.Status, req.AssigneeId, offset, limit)
	if err != nil {
		return nil, err
	}
	var protoTasks []*pb.Task
	for _, t := range tasks {
		protoTasks = append(protoTasks, toProtoTask(&t))
	}
	return &pb.ListTasksResponse{Tasks: protoTasks, Total: int32(len(protoTasks))}, nil
}

func (s *TaskServer) ChangeStatus(ctx context.Context, req *pb.ChangeStatusRequest) (*pb.ChangeStatusResponse, error) {
	userID, role, err := GetAuthContext(ctx, s.JwtService)
	if err != nil {
		return &pb.ChangeStatusResponse{Success: false}, GRPCError("unauthorized", codes.Unauthenticated)
	}
	task, err := s.Repo.GetTaskByID(req.TaskId)
	if err != nil || task == nil {
		return &pb.ChangeStatusResponse{Success: false}, GRPCError("task not found", codes.NotFound)
	}
	// Только исполнитель, создатель или админ может менять статус
	if userID != task.AssigneeID.String() && userID != task.CreatorID.String() && role != "admin" {
		return &pb.ChangeStatusResponse{Success: false}, GRPCError("forbidden", codes.PermissionDenied)
	}
	if err := s.Repo.ChangeStatus(req.TaskId, req.Status); err != nil {
		return &pb.ChangeStatusResponse{Success: false}, GRPCError("internal error", codes.Internal)
	}
	return &pb.ChangeStatusResponse{Success: true}, nil
}

func toProtoTask(t *model.Task) *pb.Task {
	var due string
	if t.DueDate != nil {
		due = t.DueDate.Format(time.RFC3339)
	}
	return &pb.Task{
		Id:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		Status:      t.Status,
		AssigneeId:  t.AssigneeID.String(),
		CreatorId:   t.CreatorID.String(),
		DueDate:     due,
		Labels:      t.Labels,
		CreatedAt:   t.CreatedAt.Format(time.RFC3339),
		UpdatedAt:   t.UpdatedAt.Format(time.RFC3339),
	}
}
