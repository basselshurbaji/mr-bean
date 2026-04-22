package auth

import "errors"

func errValidation(msg string) error {
	return errors.New(msg)
}
