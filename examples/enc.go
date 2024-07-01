package main

import (
	"fmt"

	"github.com/oarkflow/pkg/enc"
)

func main() {
	key := []byte("myverystrongpasswordo32bitlength")
	originalText := []byte("Lorem akajhsd kajhsdkasjhdkjashdk jashdkjhasdgkjhadkjahgdkjahgdkjahsgdkjhasgdkjhasg kjdasdkjhagsdkja gdkjashgd jkashgd jas")
	// Encrypt in Go
	encrypted, err := enc.EncryptCBC(string(originalText), string(key))
	if err != nil {
		fmt.Println("Encryption error:", err)
		return
	}
	fmt.Println("Encrypted:", encrypted)

	// Decrypt in Go (encrypted string from JavaScript can be used here)
	decrypted, err := enc.DecryptCFB("e8c8131d7282c5bb93036d34dc914b13:Mrpyq5u3pKBpUMtrn+9A1sZTmc5k7Flyth3O+r8Sd2s=", key)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return
	}
	fmt.Println("Decrypted:", string(decrypted))
}
