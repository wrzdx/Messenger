package core_test_utils

import (
	"errors"
	"time"
)

var (
	ID          = 1
	CreatedAt   = time.Date(2025, 1, 1, 12, 0, 0, 0, time.UTC)
	HasherError = errors.New("hash failed")
	RepoError   = errors.New("db error")
)
