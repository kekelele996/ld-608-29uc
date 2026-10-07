package utils

import (
	"fmt"

	"groundTurn/src/constants"
)

// AppError 业务错误，携带错误码与可展示消息。service 与 controller 分别包装。
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
}

func (e *AppError) Error() string { return e.Code + ": " + e.Message }

func NewAppError(status int, code, message string) *AppError {
	return &AppError{Code: code, Message: message, Status: status}
}

func ValidationError(message string) *AppError {
	return NewAppError(400, constants.ValidationFailed, message)
}

func EnumError(field, value string) *AppError {
	return NewAppError(400, constants.InvalidEnum, fmt.Sprintf(constants.InvalidEnumMessage, field, value))
}

func NotFoundError(target string, id interface{}) *AppError {
	return NewAppError(404, constants.RecordNotFound, fmt.Sprintf(constants.RecordNotFoundMessage, target, id))
}

func ForbiddenError(role, action string) *AppError {
	return NewAppError(403, constants.Forbidden, fmt.Sprintf(constants.ForbiddenMessage, role, action))
}
