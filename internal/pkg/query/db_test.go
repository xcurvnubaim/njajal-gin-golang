package query

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Test model
type Product struct {
	gorm.Model
	Name  string
	Price float64
}

// Setup mock database
func mockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	assert.NoError(t, err)

	return gormDB, mock
}

func TestApplyQuery_Search(t *testing.T) {
	db, mock := mockDB(t)
	qp := QueryParams{
		Search:        stringPtr("apple"),
		SearchColumns: []string{"name", "description"},
	}

	expectedSQL := `SELECT \* FROM "products" WHERE \(name ILIKE \$1 OR description ILIKE \$2\)`
	mock.ExpectQuery(expectedSQL).
		WithArgs("%apple%", "%apple%").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplyQuery_Filters(t *testing.T) {
	db, mock := mockDB(t)

	// Enable logging to print SQL
	db.Logger = logger.Default.LogMode(logger.Info)

	qp := QueryParams{
		Filters: &map[string]string{"category": "electronics", "price": "100"},
	}

	expectedSQL := `SELECT \* FROM "products" WHERE .+?\$1 AND .+?\$2 AND "products"."deleted_at" IS NULL`
	mock.ExpectQuery(expectedSQL).
		WithArgs("electronics", "100").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})

	// Verify all expectations
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplyQuery_Ordering(t *testing.T) {
	t.Run("ASC order", func(t *testing.T) {
		db, mock := mockDB(t)
		qp := QueryParams{OrderBy: "price", OrderDir: "asc"}

		mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"."deleted_at" IS NULL ORDER BY price ASC`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("DESC order", func(t *testing.T) {
		db, mock := mockDB(t)
		qp := QueryParams{OrderBy: "created_at", OrderDir: "DESC"}

		mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"."deleted_at" IS NULL ORDER BY created_at DESC`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestApplyQuery_Pagination(t *testing.T) {
	t.Run("Page 1 with size 10", func(t *testing.T) {
		db, mock := mockDB(t)
		qp := QueryParams{Page: 1, PageSize: 10}

		mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"."deleted_at" IS NULL LIMIT \$1`).
			WithArgs(10). // Specify the argument for the LIMIT placeholder
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Page 3 with size 20", func(t *testing.T) {
		db, mock := mockDB(t)
		qp := QueryParams{Page: 3, PageSize: 20}

		mock.ExpectQuery(`SELECT \* FROM "products" WHERE "products"."deleted_at" IS NULL LIMIT \$1 OFFSET \$2`).
			WithArgs(20, 40). // Specify the arguments for the LIMIT and OFFSET placeholders
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestApplyQuery_Combined(t *testing.T) {
	db, mock := mockDB(t)
	qp := QueryParams{
		Search:        stringPtr("laptop"),
		SearchColumns: []string{"name", "description"},
		Filters:       &map[string]string{"category": "electronics", "in_stock": "true"},
		OrderBy:       "price",
		OrderDir:      "desc",
		Page:          2,
		PageSize:      25,
	}

	expectedSQL := `SELECT \* FROM "products" WHERE \(name ILIKE \$1 OR description ILIKE \$2\) ` +
		`AND category = \$3 AND in_stock = \$4 AND "products"."deleted_at" IS NULL ` +
		`ORDER BY price DESC LIMIT \$5 OFFSET \$6`

	mock.ExpectQuery(expectedSQL).
		WithArgs("%laptop%", "%laptop%", "electronics", "true", 25, 25). // Include LIMIT and OFFSET as arguments
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

	db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestApplyQuery_EdgeCases(t *testing.T) {
	t.Run("No parameters", func(t *testing.T) {
		db, mock := mockDB(t)
		qp := QueryParams{}

		mock.ExpectQuery(`SELECT \* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Empty search columns", func(t *testing.T) {
		db, mock := mockDB(t)
		qp := QueryParams{Search: stringPtr("test")}

		mock.ExpectQuery(`SELECT \* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Zero page size", func(t *testing.T) {
		db, mock := mockDB(t)
		qp := QueryParams{PageSize: 0}

		mock.ExpectQuery(`SELECT \* FROM "products"`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}))

		db.Model(&Product{}).Scopes(qp.ApplyQuery).Find(&Product{})
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func stringPtr(s string) *string {
	return &s
}
