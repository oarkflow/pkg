package enc_test

import (
	"fmt"
	"testing"

	"github.com/oarkflow/pkg/enc"
)

func BenchmarkEncryptCFB(b *testing.B) {
	key := []byte("myverystrongpasswordo32bitlength")
	originalText := []byte("Lorem akajhsd kajhsdkasjhdkjashdk jashdkjhasdgkjhadkjahgdkjahgdkjahsgdkjhasgdkjhasg kjdasdkjhagsdkja gdkjashgd jkashgd jas")
	for i := 0; i < b.N; i++ {
		// Encrypt in Go
		_, err := enc.EncryptCFB(originalText, key)
		if err != nil {
			fmt.Println("Encryption error:", err)
			return
		}
	}
}

func BenchmarkDecryptCFB(b *testing.B) {
	key := []byte("myverystrongpasswordo32bitlength")
	encrypted := "6c77dd6a31374ac124cf203eba09ea71:LFeH76n3/ZuYwE1sKPIwMPYbH51Yro19pC8oZ3/UQzSLCLelrnMX4Zj5Sa3vcxlNTiqVqmMqrymzpAyTQSHkbMF5+krASLJ6SW5DEuP0HDBOTj4gN+8XULvt2cFSChacqjJ5aUepvHPAsinLoiKyis5jHUtR/FytUqY="
	for i := 0; i < b.N; i++ {
		_, err := enc.DecryptCFB(encrypted, key)
		if err != nil {
			fmt.Println("Encryption error:", err)
			return
		}
	}
}

func BenchmarkEncryptCBC(b *testing.B) {
	key := "myverystrongpasswordo32bitlength"
	originalText := "Lorem akajhsd kajhsdkasjhdkjashdk jashdkjhasdgkjhadkjahgdkjahgdkjahsgdkjhasgdkjhasg kjdasdkjhagsdkja gdkjashgd jkashgd jas"
	for i := 0; i < b.N; i++ {
		// Encrypt in Go
		_, err := enc.EncryptCBC(originalText, key)
		if err != nil {
			fmt.Println("Encryption error:", err)
			return
		}
	}
}

func BenchmarkDecryptCBC(b *testing.B) {
	key := "myverystrongpasswordo32bitlength"
	encrypted := "47cd83f805f2431c89f15c51ecb9ec28de1d7edf9f7db06c946656684f3c32bd96800b3d00456f302b790110266618c8ae92cb9c593dba479351ab52bdb490a7132eb90fbbd1a8cd531965e8eab05d7dc4ffd380b35d7d6eee693d3a5350bd47f78a98776dd3f3679d3f11298046e970cb4609e74b32237df63ae845b73781189159f262a9a3fdf74053452eb9fac236"
	for i := 0; i < b.N; i++ {
		_, err := enc.DecryptCBC(encrypted, key)
		if err != nil {
			fmt.Println("Encryption error:", err)
			return
		}
	}
}
