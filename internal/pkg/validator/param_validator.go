package CustomValidator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type ParamValidator struct {
	AllowedOrderByColumns []string
	AllowedFilterKeys     []string
	MaxSearchLength       int
	MaxFilterValueLength  int
	MaxPageSize           int
}

func NewParamValidator(
	allowedOrderByColumns,
	allowedFilterKeys []string,
	maxSearchLength,
	maxFilterValueLength,
	maxPageSize int,
) ParamValidator {
	return ParamValidator{
		AllowedOrderByColumns: allowedOrderByColumns,
		AllowedFilterKeys:     allowedFilterKeys,
		MaxSearchLength:       maxSearchLength,
		MaxFilterValueLength:  maxFilterValueLength,
		MaxPageSize:           maxPageSize,
	}
}

func ValidateSearch(search *string, maxlength int) error {
	if search == nil {
		return nil
	}

	if len(*search) > maxlength {
		return errors.New("Search query is too long")
	}

	if !regexp.MustCompile(`^[a-zA-Z0-9\s]*$`).MatchString(*search) {
		return errors.New("Search query contains invalid characters")
	}

	return nil
}

// ValidateFilters ensures filters contain allowed keys and safe values
func ValidateFilters(filters *map[string]string, allowedKeys []string, maxLength int) error {
	if filters == nil {
		return nil
	}

	for key, value := range *filters {
		if !contains(allowedKeys, key) {
			return fmt.Errorf("invalid filter key: %s", key)
		}
		if len(value) > maxLength {
			return fmt.Errorf("filter value for '%s' exceeds maximum length of %d", key, maxLength)
		}
	}
	return nil
}

// ValidateOrderBy ensures the column name is allowed
func ValidateOrderBy(orderBy string, allowedColumns []string) error {
	if orderBy != "" && !contains(allowedColumns, orderBy) {
		return fmt.Errorf("invalid order by column: %s", orderBy)
	}
	return nil
}

// ValidateOrderDir ensures the direction is either ASC or DESC
func ValidateOrderDir(orderDir string) error {
	orderDir = strings.ToUpper(orderDir)
	if orderDir != "ASC" && orderDir != "DESC" {
		return errors.New("invalid order direction: must be 'ASC' or 'DESC'")
	}
	return nil
}

// ValidatePage ensures the page number is greater than 0
func ValidatePage(page, pageSize, maxPageSize int) error {
	if page < 1 {
		return errors.New("page must be greater than 0")
	}

	if pageSize < 1 {
		return errors.New("page size must be greater than 0")
	}

	if pageSize > maxPageSize {
		return fmt.Errorf("page size must be less than or equal to %d", maxPageSize)
	}

	return nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
