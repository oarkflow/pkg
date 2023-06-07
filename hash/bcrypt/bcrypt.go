package bcrypt

import (
	crypt "golang.org/x/crypto/bcrypt"

	"github.com/oarkflow/pkg/str"
)

func CreateHash(password string) (string, error) {
	hash, err := crypt.GenerateFromPassword(str.ToByte(password), 8)
	return str.FromByte(hash), err
}

func ComparePasswordAndHash(password, hash string) error {
	return crypt.CompareHashAndPassword(str.ToByte(hash), str.ToByte(password))
}
