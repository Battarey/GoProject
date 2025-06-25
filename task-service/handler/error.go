package handler

import (
	"encoding/json"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func GRPCError(msg string, code codes.Code) error {
	return status.Error(code, msg)
}

func WriteJSONError(w interface{ Write([]byte) (int, error) }, code int, msg string) {
	resp, _ := json.Marshal(ErrorResponse{Error: msg})
	w.Write(resp)
}
