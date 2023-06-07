package sha512

import (
	md "crypto/sha512"
	"encoding/hex"

	"github.com/oarkflow/pkg/str"
)

func CreateHash(password string) string {
	hash := md.Sum512(str.ToByte(password))
	return hex.EncodeToString(hash[:])
}

func ComparePasswordAndHash(password, hash string) bool {
	pHash := CreateHash(password)
	return str.EqualFold(str.ToUpper(pHash), str.ToUpper(hash))
}
