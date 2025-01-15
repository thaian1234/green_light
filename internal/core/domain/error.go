package domain

import "errors"

var (
	ErrDataNotFound       = errors.New("data not found")
	ErrInternalServer     = errors.New("internal server error")
	ErrorValidation       = errors.New("validation error")
	ErrUnauthorized       = errors.New("user is unauthorized to access the resource")
	ErrForbidden          = errors.New("user is forbidden to access the resource")
	ErrorNotFound         = errors.New("not found")
	ErrNoUpdatedData      = errors.New("no data to update")
	ErrConflictingData    = errors.New("data conflicts with existing data in unique column")
	ErrTokenDuration      = errors.New("invalid token duration format")
	ErrTokenCreation      = errors.New("error creating token")
	ErrExpiredToken       = errors.New("access token has expired")
	ErrInvalidToken       = errors.New("access token is invalid")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
