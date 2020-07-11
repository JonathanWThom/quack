package secure

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"os"
)

const (
	setQuackwordError    = "Please set QUACKWORD environment variable with `export QUACKWORD=securepassword`."
	unableToDecryptError = "Failed to retrieve entries. Make sure your QUACKWORD environment variable is correct."
	unableToEncryptError = "Entry failed to save."
)

func getQuackword() (string, error) {
	quackword := os.Getenv("QUACKWORD")
	if quackword == "" {
		return "", errors.New(setQuackwordError)
	}

	return quackword, nil
}

// Decrypt reads a previously encrypted entry
func Decrypt(msg string) (string, error) {
	quackword, err := getQuackword()
	if err != nil {
		return "", err
	}

	decrypted, err := decrypt(msg, quackword)
	if err != nil {
		return "", errors.New(unableToDecryptError)
	}

	return decrypted, nil
}

// Encrypt encrypts an entry with the quackword
func Encrypt(msg string) (string, error) {
	quackword, err := getQuackword()
	if err != nil {
		return "", err
	}

	encrypted, err := encrypt(msg, quackword)
	if err != nil {
		return "", errors.New(unableToEncryptError)
	}

	return encrypted, nil
}

// EncryptWithNewQuackword encrypts an entry with a passed in quackword
func EncryptWithNewQuackword(msg string, quackword string) (string, error) {
	encrypted, err := encrypt(msg, quackword)
	if err != nil {
		return "", errors.New(unableToEncryptError)
	}

	return encrypted, nil
}

func decrypt(data, quackword string) (string, error) {
	decoded := decodeBase64(data)
	hash, err := createHash(quackword)
	if err != nil {
		return "", err
	}
	key := []byte(hash)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := decoded[:nonceSize], decoded[nonceSize:]
	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func encrypt(msg, quackword string) (string, error) {
	hash, err := createHash(quackword)
	if err != nil {
		return "", err
	}
	key := []byte(hash)
	block, _ := aes.NewCipher(key)

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(msg), nil)
	cipherString := encodeBase64(ciphertext)

	return cipherString, nil
}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func createHash(key string) (string, error) {
	hasher := md5.New()
	_, err := hasher.Write([]byte(key))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}
