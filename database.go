package main

import (
	"sync"
	"time"
)

type Database struct {
	store map[string]*Key
	mu    sync.RWMutex
}


type Key struct {
	V     string
	Exp   time.Time
}

func NewDatabase() *Database {
	return &Database{
		store: make(map[string]*Key),
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

func (db *Database) Set(key string, value string) {
	db.store[key] = &Key{V:value}

}
/*
func (db *Database) tryExpire(k string, i *Item, state *AppState) bool {
	if i.shouldExpire() {
		db.mu.Lock()
		db.Delete(k)
		db.mu.Unlock()
		return true
	}
	return false
}
*/

var db = NewDatabase()
