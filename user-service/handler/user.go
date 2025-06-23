package handler

import (
	"context"
	"errors"
	pb "user-service/proto"
)

func (s *UserServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	user, err := s.Repo.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return &pb.GetProfileResponse{
		UserId:   user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	if err := ValidateUpdateInput(req); err != nil {
		return nil, err
	}
	user, err := s.Repo.GetUserByID(req.UserId)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	user.Username = req.Username
	user.Email = req.Email
	user.Role = req.Role
	if err := s.Repo.UpdateUser(user); err != nil {
		return nil, err
	}
	return &pb.UpdateUserResponse{UserId: user.ID.String()}, nil
}

func (s *UserServer) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := s.Repo.DeleteUser(req.UserId)
	if err != nil {
		return &pb.DeleteUserResponse{Success: false}, err
	}
	return &pb.DeleteUserResponse{Success: true}, nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	offset := (int(req.Page) - 1) * int(req.PageSize)
	limit := int(req.PageSize)
	users, err := s.Repo.ListUsers(offset, limit)
	if err != nil {
		return nil, err
	}
	var userInfos []*pb.UserInfo
	for _, u := range users {
		userInfos = append(userInfos, &pb.UserInfo{
			UserId:   u.ID.String(),
			Username: u.Username,
			Email:    u.Email,
			Role:     u.Role,
		})
	}
	return &pb.ListUsersResponse{Users: userInfos, Total: int32(len(userInfos))}, nil
}
