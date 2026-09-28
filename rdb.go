package main

import (
	"encoding/gob"
	"log"
	"os"
	"path"
	"time"
)

type SnapshotTracker struct {
	keys   int
	ticker time.Ticker
	rdb    *RDBSnapshot
}

func NewSnapshotTracker(rdb *RDBSnapshot) *SnapshotTracker {
	return &SnapshotTracker{
		keys:   0,
		ticker: *time.NewTicker(time.Second * time.Duration(rdb.Secs)),
		rdb:    rdb,
	}
}

var trackers = []*SnapshotTracker{}

func InitRDBTrackers(config *Config) {
	for _, rdb := range config.rdb {
		tracker := NewSnapshotTracker(&rdb)
		trackers = append(trackers, tracker)

		go func() {
			defer tracker.ticker.Stop()

			for range tracker.ticker.C {
				log.Printf("keys changed: %d - keys required to change: %d", tracker.keys, tracker.rdb.KeysChanged)
				if tracker.keys >= tracker.rdb.KeysChanged {
					SaveRDB(config)
				}
				tracker.keys = 0
			}
		}()
	}
}

func IncrRDBTrackers() {
	for _, t := range trackers {
		t.keys++
	}
}

func SaveRDB(config *Config) {
	fp := path.Join(config.dir, config.rdbFn)
	f, err := os.OpenFile(fp, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644) // owner (read-write), everyone else (read)
	if err != nil {
		log.Println("error opening rdb file: ", err)
		return
	}
	defer f.Close()
	err = gob.NewEncoder(f).Encode(db.store)
	if err != nil {
		log.Println("error encoding rdb file: ", err)
		return
	}

	log.Println("saved RDB file")

}

func SyncRDB(conf *Config) {
	fp := path.Join(conf.dir, conf.rdbFn)
	f, err := os.Open(fp)
	if err != nil {
		log.Println("error opening rdb file: ", err)
		f.Close()
		return
	}
	defer f.Close()

	err = gob.NewDecoder(f).Decode(&db.store)
	if err != nil {
		log.Println("error decoding rdb file: ", err)
		return
	}
}

