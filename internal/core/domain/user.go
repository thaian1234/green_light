package domain

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	Plaintext *string
	Hash      []byte
}

type User struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  Password  `json:"-"`
	Activated bool      `json:"activated"`
	Version   int       `json:"-"`
}

// The Set() method calculates the bcrypt hash of a plaintext password, and stores both
// the hash and the plaintext versions in the struct.
func (p *Password) Set(plaintextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plaintextPassword), 12)
	if err != nil {
		return err
	}
	p.Plaintext = &plaintextPassword
	p.Hash = hash
	return nil
}

// The Matches() method checks whether the provided plaintext password matches the
// hashed password stored in the struct, returning true if it matches and false
// otherwise.
func (p *Password) Matches(plaintextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.Hash, []byte(plaintextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err
		}
	}
	return true, nil
}

func (u *User) PasswordMatches(plaintextPassword string) (bool, error) {
	match, err := u.Password.Matches(plaintextPassword)
	if err != nil {
		return false, err
	}
	return match, nil
}

func (u *User) IsActivated() bool {
	return u.Activated
}

func (u *User) ValidatePasswordPlaintext(password string) error {
	if password == "" {
		return errors.New("password must be provided")
	}
	if len(password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if len(password) > 72 {
		return errors.New("password must not be more than 72 characters long")
	}
	return nil
}

func (u *User) ValidateUser() error {
	if u.Name == "" {
		return errors.New("name must be provided")
	}
	if u.Email == "" {
		return errors.New("email must be provided")
	}
	if u.Password.Plaintext == nil {
		return errors.New("password must be provided")
	}
	return nil
}
