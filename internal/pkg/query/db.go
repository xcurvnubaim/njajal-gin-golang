package query

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func (qp *QueryParams) ApplyQuery(db *gorm.DB) *gorm.DB {
	// Apply search if provided
	if qp.Search != nil && len(qp.SearchColumns) > 0 {
		searchTerm := "%" + *qp.Search + "%"
		orConditions := make([]string, len(qp.SearchColumns))
		args := make([]interface{}, len(qp.SearchColumns))
		for i, col := range qp.SearchColumns {
			orConditions[i] = fmt.Sprintf("%s ILIKE ?", col)
			args[i] = searchTerm
		}
		db = db.Where(strings.Join(orConditions, " OR "), args...)
	}

	// Apply filters if provided
	if qp.Filters != nil {
		for key, value := range *qp.Filters {
			db = db.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	// Apply sorting
	if qp.OrderBy != "" {
		order := qp.OrderBy
		if strings.ToUpper(qp.OrderDir) == "DESC" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		db = db.Order(order)
	}

	// Apply pagination
	if qp.PageSize > 0 {
		db = db.Limit(qp.PageSize)
		if qp.Page > 0 {
			offset := (qp.Page - 1) * qp.PageSize
			db = db.Offset(offset)
		}
	}

	return db
}

func (qp *QueryParams) ApplyCountQuery(db *gorm.DB) *gorm.DB {
	// Apply search if provided
	if qp.Search != nil && len(qp.SearchColumns) > 0 {
		searchTerm := "%" + *qp.Search + "%"
		orConditions := make([]string, len(qp.SearchColumns))
		args := make([]interface{}, len(qp.SearchColumns))
		for i, col := range qp.SearchColumns {
			orConditions[i] = fmt.Sprintf("%s ILIKE ?", col)
			args[i] = searchTerm
		}
		db = db.Where(strings.Join(orConditions, " OR "), args...)
	}

	// Apply filters if provided
	if qp.Filters != nil {
		for key, value := range *qp.Filters {
			db = db.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}

	return db
}
