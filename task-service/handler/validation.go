package handler

import (
	"errors"
	"strings"
)

func ValidateCreateTaskInput(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title is required")
	}
	return nil
}

func ValidateUpdateTaskInput(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title is required")
	}
	return nil
}
