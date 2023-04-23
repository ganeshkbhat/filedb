package jsdbcrypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"fmt"
)

func EncryptString(str, password string) (string, error) {
	// Create a new AES cipher block using the hashed password
	hashedPassword := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(hashedPassword[:])
	if err != nil {
		return "", err
	}

	// Generate a random IV (initialization vector)
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	// Encrypt the plaintext
	plaintext := []byte(str)
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(plaintext, plaintext)

	// Combine the IV and the ciphertext into a single byte slice
	ciphertext := append(iv, plaintext...)

	// Encode the ciphertext as a base64 string and return it
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func DecryptString(hash, password string) (string, error) {
	// Decode the hash from base64
	ciphertext, err := base64.URLEncoding.DecodeString(hash)
	if err != nil {
		return "", err
	}

	// Extract the IV and the ciphertext from the combined byte slice
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Create a new AES cipher block using the hashed password
	hashedPassword := sha256.Sum256([]byte(password))
	block, err := aes.NewCipher(hashedPassword[:])
	if err != nil {
		return "", err
	}

	// Decrypt the ciphertext
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	// Convert the plaintext back to a string and return it
	return string(ciphertext), nil
}


func usage(plaintext string, password string) (string, string, error) {
	plaintext = "Hello, world!"
	password = "mysecretpassword"

	// // The plaintext string to be encrypted and decrypted
	// plaintext := "Hello, world!"
	// // The password to use for encryption and decryption
	// password := "mysecretpassword"

	// Encrypt the plaintext string
	hash, err := EncryptString(plaintext, password)
	if err != nil {
		fmt.Println("Encryption error:", err)
		return 
	}
	fmt.Println("Encrypted hash:", hash)

	// Decrypt the hash back to the plaintext string
	decryptedPlaintext, err := DecryptString(hash, password)
	if err != nil {
		fmt.Println("Decryption error:", err)
		return 
	}
	fmt.Println("Decrypted plaintext:", decryptedPlaintext)
	return hash, decryptedPlaintext, nil
}
