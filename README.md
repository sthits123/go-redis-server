# go-redis-server

A Redis-compatible in-memory data store written in Go. It speaks the [RESP protocol](https://redis.io/docs/reference/protocol-spec/) over a plain TCP socket and includes persistence, auth, memory-based eviction, transactions, and monitoring.

## Features

- **Concurrent clients** — one goroutine per connection
- **Persistence** — AOF (append-only file) and RDB (snapshot) with auto-save
- **Authentication** — `requirepass` + `AUTH`
- **Key expiration** — `EXPIRE` / `TTL`
- **Memory eviction** — LRU / LFU / random policies
- **Transactions** — `MULTI` / `EXEC` / `DISCARD`
- **Monitoring** — `MONITOR` and `INFO`
- **Configurable** via a `redis.conf`-style file

## Supported commands

```
GET   SET   DEL   EXISTS   KEYS   EXPIRE   TTL   DBSIZE
SAVE  BGSAVE  BGREWRITEAOF  FLUSHDB  INFO  MONITOR
MULTI  EXEC  DISCARD  AUTH
```

## Getting started

```sh
go build -o go-redis-server .
./go-redis-server    # listens on :6379
```

Then connect with any Redis client:

```sh
$ redis-cli
127.0.0.1:6379> SET name world
OK
127.0.0.1:6379> GET name
"world"
```

## Configuration

The server reads `./redis.conf` on startup. Missing or malformed options fall back to defaults.

```conf
dir ./data

# AOF persistence
appendonly yes
appendfilename backup.aof
appendfsync always

# RDB snapshot: save after 5 seconds or 3 key changes
save 5 3
dbfilename backup.rdb

# Memory
maxmemory 256
maxmemory-policy allkeys-lfu

# Security
requirepass mysecret
```

### Key options

| Option              | Description                                         | Default      |
| ------------------- | --------------------------------------------------- | ------------ |
| `dir`               | Directory for persistence files                     | current dir  |
| `appendonly`        | `yes` / `no` — enable the AOF file                  | `no`         |
| `appendfsync`       | `always` / `everysec` / `no` — AOF flush policy     | —            |
| `save <s> <keys>`   | Auto-snapshot rule (repeatable)                     | none         |
| `dbfilename`        | RDB snapshot filename                               | —            |
| `maxmemory`         | Memory cap in bytes (`kb`/`mb`/`gb` accepted). `0` = unlimited | `0` |
| `maxmemory-policy`  | Eviction policy: `noeviction`, `allkeys-*`, `volatile-*` | `noeviction` |
| `requirepass`       | Password required to authenticate                   | off          |