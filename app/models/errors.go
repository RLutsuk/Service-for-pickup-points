package models

import "github.com/pkg/errors"

var (
	ErrBadData              = errors.New("bad data")
	ErrBadEmail             = errors.New("invalid email format")
	ErrUserNotFound         = errors.New("login or password does not exist")
	ErrUnAuthorizated       = errors.New("user is not not authorized")
	ErrBadAuthorizated      = errors.New("user authorization error")
	ErrPickupPointDontExist = errors.New("pickup point don't exist")
	ErrInternalServer       = errors.New("internal server error")
	ErrUserExist            = errors.New("user with this email already exists ")
	ErrNotClosedReception   = errors.New("there is already an open reception in this pickup point")
	ErrNotOpenReception     = errors.New("there is no reception started")
	ErrEmptyReception       = errors.New("there are no products in this reception")
	ErrAccessDenied         = errors.New("access denied")
)
