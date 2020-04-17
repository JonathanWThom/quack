package secure

import (
	"errors"
	"fmt"
	"os"
)

const (
	setQuackwordError = "Please set QUACKWORD environment variable with `export QUACKWORD=securepassword`"
)

func getQuackword() (string, error) {
	quackword := os.Getenv("QUACKWORD")
	if quackword == "" {
		return "", errors.New(setQuackwordError)
	}

	return quackword, nil
}

func Encrypt(msg string) (string, error) {
	// check for quackword existence
	quackword, err := getQuackword()
	if err != nil {
		return "", err
	}
	fmt.Println(quackword)

	// encrypt with quackword here
	return msg, nil
}

func Decrypt(msg string) (string, error) {
	// check for quackward existence
	quackword, err := getQuackword()
	if err != nil {
		return "", err
	}
	fmt.Println(quackword)

	// decrypt with quackward here
	return msg, nil
}
