package database

import "github.com/lib/pq"

// IsUniqueConstraintError return if the passed db error is a Unique constrain violation
func IsUniqueConstraintError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
