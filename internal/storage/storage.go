package storage

import "errors"

var (
	ErrURLNotFound = errors.New("url not found")
	ErrInvalidURL  = errors.New("invalid url")
	ErrURLExists   = errors.New("url already exists")
)
