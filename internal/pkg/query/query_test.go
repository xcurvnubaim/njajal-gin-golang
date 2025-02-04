package query

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	CustomValidator "github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/validator"
)

func TestQueryParams_Parse(t *testing.T) {
	c := &gin.Context{
		Request: &http.Request{
			URL: &url.URL{
				RawQuery: "search=test&filter1=value1&filter2=value2&order_by=name&order_dir=desc&page=2&page_size=20",
			},
		},
	}

	qp := &QueryParams{}
	qp.Parse(c, "10")

	assert.Equal(t, "test", *qp.Search)
	assert.Equal(t, map[string]string{"filter1": "value1", "filter2": "value2"}, *qp.Filters)
	assert.Equal(t, "name", qp.OrderBy)
	assert.Equal(t, "desc", qp.OrderDir)
	assert.Equal(t, 2, qp.Page)
	assert.Equal(t, 20, qp.PageSize)
}

func TestQueryParams_ParseAndValidate(t *testing.T) {
	validator := CustomValidator.NewParamValidator(
		[]string{"name", "created_at"},
		[]string{"status", "category"},
		50,  // Max search length
		20,  // Max filter value length
		100, // Max page size
	)

	tests := []struct {
		name       string
		query      string
		wantErrors []string
		expected   QueryParams
	}{
		// Valid Cases
		{
			name:  "Valid full query",
			query: "search=apple&status=active&order_by=name&order_dir=desc&page=2&page_size=50",
			expected: QueryParams{
				Search:   strPtr("apple"),
				Filters:  &map[string]string{"status": "active"},
				OrderBy:  "name",
				OrderDir: "desc",
				Page:     2,
				PageSize: 50,
			},
		},
		{
			name:  "Valid minimal query",
			query: "page=1",
			expected: QueryParams{
				OrderBy:  "created_at",
				OrderDir: "asc",
				Page:     1,
				PageSize: 25, // From defaultPageSize
			},
		},

		// Invalid Cases
		{
			name:       "Invalid search characters",
			query:      "search=apple!",
			wantErrors: []string{"Search query contains invalid characters"},
		},
		{
			name:       "Oversized search",
			query:      "search=" + strings.Repeat("a", 51),
			wantErrors: []string{"Search query is too long"},
		},
		{
			name:       "Invalid filter key",
			query:      "invalid_key=value",
			wantErrors: []string{"invalid filter key: invalid_key"},
		},
		{
			name:       "Oversized filter value",
			query:      "status=" + strings.Repeat("a", 21),
			wantErrors: []string{"exceeds maximum length of 20"},
		},
		{
			name:       "Invalid order column",
			query:      "order_by=invalid_column",
			wantErrors: []string{"invalid order by column: invalid_column"},
		},
		{
			name:       "Invalid order direction",
			query:      "order_dir=up",
			wantErrors: []string{"invalid order direction: must be 'ASC' or 'DESC'"},
		},
		{
			name:  "Invalid page values",
			query: "page=0&page_size=101",
			wantErrors: []string{
				"page must be greater than 0",
			},
		},
		{
			name:  "Invalid page values",
			query: "page=1&page_size=101",
			wantErrors: []string{
				"page size must be less than or equal to 100",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup Gin context with test query
			c := createTestContext(tt.query)

			// Parse parameters
			qp := &QueryParams{}
			qp.Parse(c, "25") // Default page size 25

			// Validate parameters
			err := qp.Validate(validator)

			// Verify expected values for valid cases
			if len(tt.wantErrors) == 0 {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected.Search, qp.Search)
				assert.Equal(t, tt.expected.Filters, qp.Filters)
				assert.Equal(t, tt.expected.OrderBy, qp.OrderBy)
				assert.Equal(t, tt.expected.OrderDir, qp.OrderDir)
				assert.Equal(t, tt.expected.Page, qp.Page)
				assert.Equal(t, tt.expected.PageSize, qp.PageSize)
			} else {
				assert.Error(t, err)
				for _, msg := range tt.wantErrors {
					assert.Contains(t, err.Error(), msg)
				}
			}
		})
	}
}

func createTestContext(query string) *gin.Context {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{
		URL: &url.URL{RawQuery: query},
	}
	return c
}

func strPtr(s string) *string {
	return &s
}
