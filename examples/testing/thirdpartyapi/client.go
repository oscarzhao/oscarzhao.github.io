package thirdpartyapi

import (
	"errors"
)

//go:generate mockery -name=Client

// Client defines operations a third party service has
type Client interface {
	Get(key string) (data interface{}, err error)
}

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("key not found")
)
