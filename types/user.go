package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost      = 12
	minFirstNameLen = 3
	minLastNameLen  = 3
	minPasswordLen  = 8
)

var (
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	lowercaseRegex   = regexp.MustCompile(`[a-z]`)
	uppercaseRegex   = regexp.MustCompile(`[A-Z]`)
	numberRegex      = regexp.MustCompile(`\d`)
	specialCharRegex = regexp.MustCompile(`[@$!%*?&#]`)
)

type UpdateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	mapped := bson.M{}
	if len(p.FirstName) > 0 {
		mapped["first_name"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		mapped["last_name"] = p.LastName
	}
	return mapped
}

type CreateUserParams struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func validateEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func validatePassword(password string) bool {
	if len(password) < minPasswordLen {
		return false
	}
	if !lowercaseRegex.MatchString(password) {
		return false
	}
	if !uppercaseRegex.MatchString(password) {
		return false
	}
	if !numberRegex.MatchString(password) {
		return false
	}
	if !specialCharRegex.MatchString(password) {
		return false
	}
	return true
}

func (params CreateUserParams) Validate() error {
	if len(params.FirstName) < minFirstNameLen {
		return fmt.Errorf("first name must be at least %d characters long", minFirstNameLen)
	}

	if len(params.LastName) < minLastNameLen {
		return fmt.Errorf("last name must be at least %d characters long", minLastNameLen)
	}

	if !validateEmail(params.Email) {
		return fmt.Errorf("invalid email address")
	}

	if !validatePassword(params.Password) {
		return fmt.Errorf("password must be at least %d characters long and contain at least one uppercase letter, one lowercase letter, and one digit", minPasswordLen)
	}

	return nil
}

// ValidatePassword and compare the password with the encrypted password
func IsPasswordValid(encryptedPw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encryptedPw), []byte(pw)) == nil
}

// User struct to represent a user
type User struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	FirstName         string             `json:"first_name" bson:"first_name"`
	LastName          string             `json:"last_name" bson:"last_name"`
	Email             string             `json:"email" bson:"email"`
	EncryptedPassword string             `json:"-" bson:"encrypted_password"`
}

func CreateNewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
