package errors

import (
	"errors"
	"fmt"
)

func New(message string) error {
	return errors.New(message)
}

func Wrap(prefix string, err error, message string, args ...any) error {
	return fmt.Errorf("[%s] %s: %w", prefix, fmt.Sprintf(message, args...), err)
}
