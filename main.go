package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	log.Println("Reading config file...")
	conf := readConf("./redis.conf")
	state := newAppState(conf)
	if conf.aofEnabled {
		log.Println("Syncing AOF...")
		state.aof.Sync()
	}
	

	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal("Cannot listen on port 6379: ", err)
	}
	defer listener.Close()
	log.Println("listening on :6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		log.Println("connection accepted")

		go func() {
			handleConn(conn, state)
		}()
	}

	
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


func handleConn(conn net.Conn, state *AppState) {
	log.Println("accepted new connection: ", conn.LocalAddr().String())
	c := NewClient(conn)
	r := bufio.NewReader(conn)

	for {
		v := Value{typ: ARRAY}
		if err := v.readArray(r); err != nil {
			log.Println(err)
			break
		}
		handle(c, &v, state)
	}
	log.Println("connection closed: ", conn.LocalAddr().String())
}