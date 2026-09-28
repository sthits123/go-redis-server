package main

import "time"

type AppState struct {
	conf              *Config
	aof               *Aof
	bgsaveRunning     bool
	aofRewriteRunning bool
	dbCopy            map[string]*Item
	tx                *Transaction
	monitors          []*Client
	serverStart       time.Time
	clientCount       int
	peakMem           int64
	info              *Info
	rdbStats          RDBStats
	aofStats          AOFStats
	generalStats      GeneralStats
}

type RDBStats struct {
	rdb_last_save_ts int64
	rdb_saves        int
}

type AOFStats struct {
	aof_rewrites int
}

type GeneralStats struct {
	total_connections_received int
	total_commands_processed   int
	expired_keys               int
	evicted_keys               int
}

