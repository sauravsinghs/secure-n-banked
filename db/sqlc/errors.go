package db

import (
	"database/sql"

	"github.com/lib/pq"
)

// Standardized PostgreSQL error codes used in the application
const (
	UniqueViolation      = "unique_violation"
	ForeignKeyViolation  = "foreign_key_violation"
)

// Helper error variables for convenience in tests and handlers
var (
	ErrRecordNotFound   = sql.ErrNoRows
	ErrUniqueViolation  = &pq.Error{Code: "23505"}
)

// ErrorCode extracts the Postgres error code name from the given error.
// Returns an empty string if the error is not a *pq.Error.
func ErrorCode(err error) string {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code.Name()
	}
	return ""
} 