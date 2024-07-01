package enc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/oarkflow/pkg/pkcs7"
)

// EncryptCBC encrypts plain text string into cipher text string
func EncryptCBC(unencrypted string, password string) (string, error) {
	key := []byte(password)
	plainText := []byte(unencrypted)
	plainText, err := pkcs7.Pad(plainText, aes.BlockSize)
	if err != nil {
		return "", fmt.Errorf(`plainText: "%s" has error`, plainText)
	}
	if len(plainText)%aes.BlockSize != 0 {
		err := fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)

	return fmt.Sprintf("%x", cipherText), nil
}

// DecryptCBC decrypts cipher text string into plain text string
func DecryptCBC(encrypted string, password string) (string, error) {
	key := []byte(password)
	cipherText, err := hex.DecodeString(encrypted)
	if err != nil {
		fmt.Printf("%v", err)
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("%v", err)
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		return "", errors.New("cipherText too short")
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		return "", errors.New("cipherText is not a multiple of the block size")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)

	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	return fmt.Sprintf("%s", cipherText), nil
}

func EncryptCFB(data []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], data)

	ivHex := hex.EncodeToString(iv)
	encryptedData := base64.StdEncoding.EncodeToString(ciphertext[aes.BlockSize:])
	return ivHex + ":" + encryptedData, nil
}

func DecryptCFB(encrypted string, key []byte) ([]byte, error) {
	parts := strings.Split(encrypted, ":")
	ivHex := parts[0]
	encryptedData := parts[1]

	iv, err := hex.DecodeString(ivHex)
	if err != nil {
		return nil, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedData)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

/*
// Encrypt function
function encrypt(data) {
    const iv = CryptoJS.lib.WordArray.random(16); // Generate a random IV
    const encrypted = CryptoJS.AES.encrypt(data, key, { iv: iv, mode: CryptoJS.mode.CFB, padding: CryptoJS.pad.NoPadding });
    const ivHex = CryptoJS.enc.Hex.stringify(iv);
    const encryptedData = encrypted.ciphertext.toString(CryptoJS.enc.Base64);
    return ivHex + ":" + encryptedData;
}

// Decrypt function
function decrypt(encrypted) {
    const parts = encrypted.split(':');
    const ivHex = parts[0];
    const encryptedData = parts[1];
    const iv = CryptoJS.enc.Hex.parse(ivHex);
    const decrypted = CryptoJS.AES.decrypt(CryptoJS.lib.CipherParams.create({
        ciphertext: CryptoJS.enc.Base64.parse(encryptedData)
    }), key, { iv: iv, mode: CryptoJS.mode.CFB, padding: CryptoJS.pad.NoPadding });
    return decrypted.toString(CryptoJS.enc.Utf8);
}
*/
