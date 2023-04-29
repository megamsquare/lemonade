package user

import (
	"context"
	"errors"
	"log"
	"time"
)

func (u *userMemoryStore) QueryByID(ctx context.Context, traceID string, id int64) (*User, error) {
	user, ok := u.Table[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (u *userMemoryStore) QueryAll(ctx context.Context, traceID string) ([]User, error) {
	values := []User{}

	if len(u.Table) < 1 {
		return nil, ErrNotFound
	}

	for _, value := range u.Table {
		values = append(values, *value)
	}

	return values, nil
}

func (u *userMemoryStore) Create(ctx context.Context, traceID string, newUser *NewUser) error {
	user := User{
		ID:                 int64(len(u.Table) + 1),
		Username:           newUser.Username,
		Balance:            newUser.Balance,
		VerificationStatus: newUser.VerificationStatus,
	}
	verifyUser := VerifyInfo{
		ID:     user.ID,
		Status: newUser.VerificationStatus,
	}

	u.Table[user.ID] = &user
	u.VerifyQueue = append(u.VerifyQueue, verifyUser)
	// u.VerifyQueue[user.ID] = &verifyUser

	return nil
}

func (u *userMemoryStore) Update(ctx context.Context, traceID string, user *User) error {
	savedUser, err := u.QueryByID(ctx, traceID, user.ID)
	if err != nil {
		return err
	}
	u.Table[savedUser.ID] = user

	return nil
}

func (u *userMemoryStore) IsVerified(ctx context.Context, traceID string, id int64) (bool, error) {
	// transaction, ok := u.VerifyQueue[id]
	for _, verifyUser := range u.VerifyQueue {
		if verifyUser.ID == id {
			if !verifyUser.Status {
				return false, nil
			} else {
				return true, nil
			}
		}
	}

	return false, ErrNotFound
}

func (u *userMemoryStore) verifyUser(traceID string, user VerifyInfo) (VerifyInfo, error) {
	if !user.Status {
		user.Status = true
		return user, nil
	}
	return VerifyInfo{}, errors.New("User is already verified")
}

func (u *userMemoryStore) VerifyProcess() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			for i, user := range u.VerifyQueue {
				log.Println("Verification in progress")
				go func(i int, user VerifyInfo) {
					if !user.Status {
						v, err := u.verifyUser("Verification", user)
						if err != nil {
							log.Println("Error verifying user:", err)
						}
						u.VerifyQueue[i] = v
						savedUser, err := u.QueryByID(context.Background(), "Verify User", v.ID)
						if err != nil {
							log.Println(err)
						}
						savedUser.VerificationStatus = v.Status
						u.Table[savedUser.ID] = savedUser
						log.Println("Verified User: ", savedUser.ID)
					}
				}(i, user)
			}
		}
	}
}
