package models

import "golang.org/x/crypto/bcrypt"

type Profile struct {
	Id              uint64  `json:"-"`
	FirstName       string  `json:"first_name"`
	LastName        string  `json:"last_name"`
	Email           string  `json:"email" gorm:"unique"`
	Password        []byte  `json:"-"`
	PasswordEntropy float64 `json:"-"` // optional
}

func (profile *Profile) SetPassword(password []byte) error {
	var err = new(error)
	profile.Password, *err = bcrypt.GenerateFromPassword(password, 16) // default is 10, range 4-31
	return *err
}

func (profile *Profile) VerifyPassword(password []byte) error {
	return bcrypt.CompareHashAndPassword(profile.Password, password)
}
