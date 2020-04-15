package cmd

import (
	"github.com/jonathanwthom/quack/storage"
)

type fakeStorage struct{}

func (s *fakeStorage) Create(msg string) error {
	return nil
}

func (s *fakeStorage) Read() ([]storage.Entry, error) {
	return []storage.Entry{}, nil
}
