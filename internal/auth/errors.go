package auth

import "errors"

var ErrInvalidCredentials = errors.New("invalid credentials")
var ErrNotAuthenticated = errors.New("not authenticated")
var ErrEmailAlreadyTaken = errors.New("email already taken")
