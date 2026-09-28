package main

import (
	"log"
	"net"
)

type Handler func(*Value) *Value

var Handlers = map[string]Handler{
	"COMMAND":      command,
	"GET":          get,
	"SET":          set,
	
}



func handle(conn net.Conn, v *Value) {
	cmd :=v.array[0].bulk
	handler, ok := Handlers[cmd]
	if !ok {
		log.Println("invalid command: ", cmd)
		return
	}

	reply:=handler(v)
	w:=NewWriter(conn)
	w.Write(reply)


}

func get(v *Value) *Value {
   args := v.array[1:]
	if len(args) != 1 {
		return &Value{typ: ERROR, err: "ERR invalid number of arguments for 'GET' command"}
	}

	name := args[0].bulk
	db.mu.RLock()
	val,ok:= db.store[name]
	db.mu.RUnlock()

	if !ok {
		return &Value{typ: NULL}
	}

	return &Value{typ: BULK, bulk:val}
}

func set(v *Value) *Value {
	args := v.array[1:]
	if len(args) != 2 {
		return &Value{typ: ERROR, err: "ERR invalid number of arguments for 'SET' command"}
	}

	
	key:= args[0].bulk
	value := args[1].bulk
	db.mu.Lock()
	db.store[key] = value
	db.mu.Unlock()

	return &Value{typ: STRING, str:  "OK"}
}


func command(v *Value) *Value {
	return &Value{typ:STRING,str:"OK"}
}