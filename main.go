package main

import (
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	log.Println("Reading config file...")
	conf := readConf("./redis.conf")
	state := newAppState(conf)
	if state.conf.aofEnabled {
		log.Println("Syncing AOF...")
		state.aof.Sync()
	}

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal("Cannot listen on port 6379: ", err)
	}
	defer listener.Close()
	log.Println("listening on :6379")

	conn, err := listener.Accept()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Println("connection accepted")

	for {
		v := Value{typ: ARRAY}
		if err := v.readArray(conn); err != nil {
			if err != io.EOF {
				log.Println("protocol error, closing connection: ", err)
			}
			return
		}
		handle(conn, &v, state)
	}
}

type AppState struct {
	conf *Config
	aof  *Aof
}

func newAppState(conf *Config) *AppState {
	state := AppState{conf: conf}

	if conf.aofEnabled {
		state.aof = NewAof(conf)

		if conf.aofFsync == EverySec {
			go func() {
				t := time.NewTicker(time.Second)
				defer t.Stop()

				for range t.C {
					state.aof.w.Flush()
				}
			}()
		}
	}

	return &state
}