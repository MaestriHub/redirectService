package pkg

import (
	"net/http"
	"strconv"
)

type ErrorS struct {
	Message string
	Status  int
}

func (e *ErrorS) Error() string {
	return strconv.Itoa(e.Status) + " " + e.Message
}

func NewErrorS(msg string, status int) *ErrorS {
	return &ErrorS{msg, status}
}

func NewBadRequestError(msg string) *ErrorS {
	return NewErrorS(msg, http.StatusBadRequest)
}

func NewNotFoundError(msg string) *ErrorS {
	return NewErrorS(msg, http.StatusNotFound)
}

func NewInternalServerError(msg string) *ErrorS {
	return NewErrorS(msg, http.StatusInternalServerError)
}

func NewUnauthorizedError(msg string) *ErrorS {
	return NewErrorS(msg, http.StatusUnauthorized)
}

func NewForbiddenError(msg string) *ErrorS {
	return NewErrorS(msg, http.StatusForbidden)
}

func NewConflictError(msg string) *ErrorS {
	return NewErrorS(msg, http.StatusConflict)
}

func NewNotImplementedError(msg string) *ErrorS {
	return NewErrorS(msg, http.StatusNotImplemented)
}
