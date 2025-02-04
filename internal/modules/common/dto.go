package common

import "github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/query"

type (
	PaginationResponseDTO[D any] struct {
		Data *D                    `json:"data"`
		Meta *query.PaginationMeta `json:"meta"`
	}
)
