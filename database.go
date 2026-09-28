package main

import (
	"sync"
)

type Database struct {
	store map[string]string
	mu    sync.RWMutex
}

func NewDatabase() *Database {
	return &Database{
		store: make(map[string]string),
		mu:    sync.RWMutex{},
	}
}

var db = NewDatabase()
