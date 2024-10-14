package password

import (
	"crypto/subtle"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHasher is an interface type to declare a general-purpose password management tool.
type PasswordHasher interface {
	HashPassword(string) (string, error)
	VerifyPassword(string, string) bool
}
type DummyPasswordHasher struct{}

type BcryptPasswordHasher struct {
	Cost int
}

var (
	_ PasswordHasher = DummyPasswordHasher{}
	_ PasswordHasher = BcryptPasswordHasher{0}
)

var preferredHashers = []PasswordHasher{
	BcryptPasswordHasher{},
}

func HashPassword(password string) (string, error) {
	return hashPasswordWithHashers(password, preferredHashers)
}

func hashPasswordWithHashers(password string, hashers []PasswordHasher) (string, error) {
	// Even though good hashers will disallow blank passwords, let's be explicit that ALL BLANK PASSWORDS ARE INVALID.  Full stop.
	if password == "" {
		return "", fmt.Errorf("blank passwords are not allowed")
	}
	return hashers[0].HashPassword(password)
}

// HashPassword creates a one-way digest ("hash") of a password.  In the case of Bcrypt, a pseudorandom salt is included automatically by the underlying library.
func (h DummyPasswordHasher) HashPassword(password string) (string, error) {
	return password, nil
}

// VerifyPassword validates whether a one-way digest ("hash") of a password was created from a given plaintext password.
func (h DummyPasswordHasher) VerifyPassword(password, hashedPassword string) bool {
	return 1 == subtle.ConstantTimeCompare([]byte(password), []byte(hashedPassword))
}

func (h BcryptPasswordHasher) HashPassword(password string) (string, error) {
	cost := h.Cost
	if cost < bcrypt.DefaultCost {
		cost = bcrypt.DefaultCost
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		hashedPassword = []byte("")
	}
	return string(hashedPassword), err
}

// VerifyPassword validates whether a one-way digest ("hash") of a password was created from a given plaintext password.
func (h BcryptPasswordHasher) VerifyPassword(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
