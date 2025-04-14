package models

import "github.com/pkg/errors"

// business logic errors
var (
	ErrPickupPointDontExist = errors.New("pickup point don't exist")
	ErrInternalServer       = errors.New("internal server error")
	ErrNotClosedReception   = errors.New("there is already an open reception in this pickup point")
	ErrNotOpenReception     = errors.New("there is no reception started")
	ErrEmptyReception       = errors.New("there are no products in this reception")
)

// data errors
var (
	ErrBadData  = errors.New("bad data")
	ErrBadEmail = errors.New("invalid email format")
)

// user erros
var (
	ErrAccessDenied    = errors.New("access denied")
	ErrUnAuthorizated  = errors.New("user is not not authorized")
	ErrBadAuthorizated = errors.New("user authorization error")
	ErrUserNotFound    = errors.New("login or password does not exist")
	ErrUserExist       = errors.New("user with this email already exists")
)

// swagger

type ErrorResponse struct {
	Message string `json:"message" example:"string error"`
}
