package errors

import (
	"net/http"
)

type Error struct {
	Code int    `json:"code"`
	Err  string `json:"errors"`
}

func NoErr() Error {
	return Error{}
}
func (e Error) Error() string {
	return e.Err
}

func NewError(code int, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}
func ErrUnAuthorized(res string) *Error {
	return &Error{Code: http.StatusUnauthorized, Err: res + "  ,unauthorized request"}

}
func ErrResourceNotFound(res string) *Error {
	return &Error{
		Code: http.StatusNotFound,
		Err:  res + "  ,resource not found",
	}
}

func ErrDB(res string) *Error {
	return &Error{
		Code: http.StatusInternalServerError,
		Err:  res + "  ,resource not found",
	}
}
func ErrServer(res string) *Error {
	return &Error{
		Code: http.StatusInternalServerError,
		Err:  res + "  ,server errors",
	}
}

func ErrBadRequest(res string) *Error {
	return &Error{
		Code: http.StatusBadRequest,
		Err:  res + "  ,invalid JSON request",
	}
}
