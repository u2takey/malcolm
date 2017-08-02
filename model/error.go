package model

import (
	"errors"
	"strings"
)

// WarpErrors wrap errors
func WarpErrors(errs ...error) error {
	msg := []string{}
	for _, err := range errs {
		if err != nil {
			msg = append(msg, err.Error())
		}
	}
	if len(msg) > 0 {
		return errors.New(strings.Join(msg, "\n"))
	}
	return nil
}
