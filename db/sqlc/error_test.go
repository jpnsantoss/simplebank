package db

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/require"
)

func TestErrorCode(t *testing.T) {
	err := errors.New("some error")
	code := ErrorCode(err)
	require.Equal(t, "", code, "expected empty error code")

	pgErr := &pgconn.PgError{
		Code: ForeignKeyViolation,
	}
	code = ErrorCode(pgErr)
	require.Equal(t, ForeignKeyViolation, code, "expected ForeignKeyViolation error code")
}
