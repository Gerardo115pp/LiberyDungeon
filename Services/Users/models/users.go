package models

import (
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUID       string `json:"uuid"`
	Username   string `json:"username"`
	SecretHash string `json:"secret_hash"`
}

type UserEntry struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
}

type UserIdentity struct {
	UUID          string   `json:"uuid"`
	Username      string   `json:"username"`
	RoleHierarchy int      `json:"role_hierarchy"`
	Grants        []string `json:"grants"`
}

type CreateUserParams struct {
	Username        string `json:"username"`
	PlainTextSecret string `json:"secret"`
}

func CreateNewUser(params CreateUserParams) (user *User, err error) {
	var user_uuid string = uuid.NewString()

	user = &User{
		UUID:     user_uuid,
		Username: params.Username,
	}

	user.UpdateSecret(params.PlainTextSecret)

	return user, nil
}

func CreateUserIdentity(user *User, role_hierarchy int, grants []string) *UserIdentity {
	return &UserIdentity{
		UUID:          user.UUID,
		Username:      user.Username,
		RoleHierarchy: role_hierarchy,
		Grants:        grants,
	}
}

func (user *User) CompareSecret(plain_text_secret string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.SecretHash), []byte(plain_text_secret))
}

func (user User) GetAsEntry() UserEntry {
	return UserEntry{
		UUID:     user.UUID,
		Username: user.Username,
	}
}

func (user *User) UpdateSecret(new_secret string) error {
	new_secret_hash, err := bcrypt.GenerateFromPassword([]byte(new_secret), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.SecretHash = string(new_secret_hash)

	return nil
}

func (user *User) UpdateFromOther(other *User) error {
	if other.UUID != user.UUID {
		return fmt.Errorf("User UUIDs do not match")
	}

	user.Username = other.Username

	if other.SecretHash != "" {
		if err := user.UpdateSecret(other.SecretHash); err != nil {
			return err
		}
	}

	return nil
}
