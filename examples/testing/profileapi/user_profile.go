package profilesdk

import (
	"errors"
)

//go:generate mockery -name=UserProfileAPI

// Profile defines fields for a user's profile
type Profile struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// UserProfileAPI defines operations allowed by User Profile service
type UserProfileAPI interface {
	Get(name string) (data *Profile, err error)
}

var (
	// ErrNotFound ...
	ErrNotFound = errors.New("user not found")
)
