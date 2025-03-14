package custom_errors

import "errors"

var (
	ErrNoRows = errors.New("not found")
)
