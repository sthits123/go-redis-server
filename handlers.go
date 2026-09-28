package main

import (
	"log"
	"net"
	"strings"
)

type Handler func(*Value, *AppState) *Value

var Handlers = map[string]Handler{
	"COMMAND": command,
	"GET":     get,
	"SET":     set,
}

func handle(conn net.Conn, v *Value, state *AppState) {
	if len(v.array) == 0 {
		log.Println("empty command received")
		NewWriter(conn).Write(&Value{typ: ERROR, err: "ERR empty command"})
		return
	}

	cmd := strings.ToUpper(v.array[0].bulk)
	handler, ok := Handlers[cmd]
	if !ok {
		log.Println("invalid command: ", cmd)
		return
	}

	reply := handler(v, state)
	w := NewWriter(conn)
	w.Write(reply)
	w.Flush()

}

func get(v *Value, state *AppState) *Value {
	args := v.array[1:]
	if len(args) != 1 {
		return &Value{typ: ERROR, err: "ERR invalid number of arguments for 'GET' command"}
	}

	name := args[0].bulk
	db.mu.RLock()
	val, ok := db.store[name]
	db.mu.RUnlock()

	if !ok {
		return &Value{typ: NULL}
	}

	return &Value{typ: BULK, bulk: val}
}

func set(v *Value, state *AppState) *Value {
	args := v.array[1:]
	if len(args) != 2 {
		return &Value{typ: ERROR, err: "ERR invalid number of arguments for 'SET' command"}
	}

	key := args[0].bulk
	value := args[1].bulk
	db.mu.Lock()
	db.store[key] = value
	if state.conf.aofEnabled {
		log.Println("writing to aof")
		state.aof.w.Write(v)
		if state.conf.aofFsync == Always {
			state.aof.w.Flush()
		}
	}
	db.mu.Unlock()

	return &Value{typ: STRING, str: "OK"}
}

func command(v *Value, state *AppState) *Value {
	return &Value{typ: STRING, str: "OK"}
}
