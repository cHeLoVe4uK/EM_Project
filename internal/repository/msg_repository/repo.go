package msg_repository

import "errors"

var (
	ErrNotAllowed      = errors.New("not allowed")
	ErrMessageNotFound = errors.New("message not found")
)
