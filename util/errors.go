package util

import (
	"net/http"
	"strings"

	"github.com/juju/errors"
)

// WrapError ...
func WrapError(err error) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "Duplicate entry"):
		return errors.AlreadyExistsf("Duplicate entry")
	}
	return err
}

// CauseError ...
func CauseError(err error, code int, msg string) (int, string) {
	switch {
	case errors.IsNotFound(err), errors.IsUserNotFound(err):
		return http.StatusNotFound, errors.Cause(err).Error()
	case errors.IsAlreadyExists(err):
		return http.StatusConflict, errors.Cause(err).Error()
	}
	return code, msg
}
