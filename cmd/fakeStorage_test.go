package cmd

import (
	"github.com/jonathanwthom/quack/storage"
)

type fakeStorage struct{}

func (s *fakeStorage) Create(msg string) error {
	return nil
}

var entriesMock []storage.Entry
var errorMock error

func (s *fakeStorage) Read() ([]storage.Entry, error) {
	return entriesMock, errorMock
}
