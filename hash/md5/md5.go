package md5

import (
	md "crypto/md5"
	"encoding/hex"

	"github.com/oarkflow/pkg/str"
)

func CreateHash(password string) string {
	hash := md.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func ComparePasswordAndHash(password, hash string) bool {
	pHash := CreateHash(password)
	return str.EqualFold(str.ToUpper(pHash), str.ToUpper(hash))
}
