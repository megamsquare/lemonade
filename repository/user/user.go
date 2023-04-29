package user

import (
	"context"
	"errors"
)

type User struct {
	ID                 int64
	Username           string
	Balance            float64
	VerificationStatus bool
}

type VerifyInfo struct {
	ID     int64
	Status bool
}

type userMemoryStore struct {
	Table       map[int64]*User
	VerifyQueue []VerifyInfo
}

type NewUser struct {
	Username           string `json:"username"`
	Balance            float64
	VerificationStatus bool
}

var (
	ErrNotFound = errors.New("user not found")
)

type UserRepository interface {
	QueryByID(ctx context.Context, traceID string, id int64) (*User, error)
	QueryAll(ctx context.Context, traceID string) ([]User, error)
	Create(ctx context.Context, traceID string, newUser *NewUser) error
	Update(ctx context.Context, traceID string, user *User) error
	IsVerified(ctx context.Context, traceID string, id int64) (bool, error)
	VerifyProcess()
}

func NewUserMemoryStore() UserRepository {
	return &userMemoryStore{
		Table:       make(map[int64]*User),
		VerifyQueue: []VerifyInfo{},
	}
}
