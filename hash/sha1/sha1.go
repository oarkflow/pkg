package sha1

import (
	sh "crypto/sha1"
	"encoding/hex"

	"github.com/oarkflow/pkg/str"
)

func CreateHash(password string) string {
	hash := sh.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func ComparePasswordAndHash(password, hash string) bool {
	pHash := CreateHash(password)
	return str.EqualFold(str.ToUpper(pHash), str.ToUpper(hash))
}
