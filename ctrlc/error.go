package ctrlc

import "context"

type Error struct {
	error
}

func AppError(err error) *Error {
	return &Error{err}
}

func (c *Error) Shutdown(ctx context.Context) error {
	return c
}
