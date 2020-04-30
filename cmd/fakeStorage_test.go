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

var readByKeyMock storage.Entry
var readByKeyErrorMock error

func (s *fakeStorage) ReadByKey(key string) (storage.Entry, error) {
	return readByKeyMock, readByKeyErrorMock
}

func (s *fakeStorage) Delete(key string) error {
	return nil
}
