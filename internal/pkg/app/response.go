package app

import "github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/query"

type SuccessResponse[D any] struct {
	Status  bool    `json:"status"`
	Message string  `json:"message"`
	Data    *D      `json:"data"`
	Error   *string `json:"error"`
}

type ErrorResponse struct {
	Status  bool    `json:"status"`
	Message string  `json:"message"`
	Data    *string `json:"data"`
	Error   *string `json:"error"`
}

type PaginationResponse[D any] struct {
	Status  bool                  `json:"status"`
	Message string                `json:"message"`
	Meta    *query.PaginationMeta `json:"meta"`
	Data    *D                    `json:"data"`
}

func NewSuccessResponse[D any](message string, data *D) *SuccessResponse[D] {
	return &SuccessResponse[D]{
		Status:  true,
		Message: message,
		Data:    data,
		Error:   nil,
	}
}

func NewErrorResponse(message string, err *string) *ErrorResponse {
	return &ErrorResponse{
		Status:  false,
		Message: message,
		Data:    nil,
		Error:   err,
	}
}

func NewPaginationResponse[D any](message string, meta *query.PaginationMeta, data *D) *PaginationResponse[D] {
	return &PaginationResponse[D]{
		Status:  true,
		Message: message,
		Meta:    meta,
		Data:    data,
	}
}