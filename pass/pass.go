package pass

import (
	"golang.org/x/crypto/bcrypt"
)

type Pass struct {
}

func New() *Password {
	return &Pass{}
}

// Хеширует строку
func (self *Pass) Hash(p string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(p), 8)
	return string(hash), err
}

// Сверяет хеш со строкой
func (self *Pass) Compare(h, p string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(h), []byte(p))
	if err != nil {
		return false
	}
	return true
}
