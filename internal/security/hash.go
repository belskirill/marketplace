package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashPassword), err
}

func CheckPasswordHash(password, hashPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))

	return err == nil
}

type Car interface {
	Drive()
}

type BMW struct {
	model string
}

func (b *BMW) Drive() {
	fmt.Println("BMW drive")
}

type Mercedes struct {
	model string
}

func (m *Mercedes) Drive() {
	fmt.Println("Mercedes drive")
}

func AssertType(cr Car) {
	if c, ok := cr.(*BMW); ok {
		fmt.Println("is it BMW:", c.model)
	}
}
