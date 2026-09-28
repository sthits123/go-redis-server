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

func (db *Database) Delete(k string) {
	_, ok := db.store[k]
	if !ok {
		return 
	}

	delete(db.store, k)
}

var db = NewDatabase()
