package main

import (
	"log"
	"path/filepath"
	"strconv"
    "time"
)

type Handler func(*Client,*Value, *AppState) *Value



var Handlers = map[string]Handler{
	"COMMAND": command,
	"GET":     get,
	"SET":     set,
	"DEL":     del,
	"EXISTS":  exists,
	"KEYS":    keys,
	"EXPIRE":  expire,
	"TTL":     ttl,
	"MULTI":   multi,
	"EXEC":     _exec,
	"DISCARD":  discard,
}

func handle(c *Client,v *Value, state *AppState) {

	cmd := v.array[0].bulk
	handler, ok := Handlers[cmd]
	w := NewWriter(c.conn)

	if !ok {
		w.Write(&Value{typ: ERROR, err: "ERR invalid command"})
		w.Flush()
		return
	}

	if state.tx != nil && cmd != "EXEC" && cmd != "DISCARD" {
		txCmd := TxCommand{v: v, handler: handler}
		state.tx.cmds = append(state.tx.cmds, &txCmd)
		w.Write(&Value{typ: STRING, str: "QUEUED"})
		w.Flush()
		return
	}

	reply := handler(c, v, state)
	w.Write(reply)
	w.Flush()


	

}

func get(c *Client,v *Value, state *AppState) *Value {
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

	return &Value{typ: BULK, bulk: val.V}
}

func set(c *Client,v *Value, state *AppState) *Value {
	args := v.array[1:]
	if len(args) != 2 {
		return &Value{typ: ERROR, err: "ERR invalid number of arguments for 'SET' command"}
	}

	key := args[0].bulk
	value := args[1].bulk
	db.mu.Lock()
	db.Set(key,value)
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

func command(c *Client,v *Value, state *AppState) *Value {
	return &Value{typ: STRING, str: "OK"}
}


func del(c *Client,v *Value, state *AppState) *Value {
	args := v.array[1:]
	var n int

	db.mu.Lock()
	for _, arg := range args {
		_, ok := db.store[arg.bulk]
		db.Delete(arg.bulk)
		if ok {
			n++
		}
	}
	db.mu.Unlock()

	return &Value{typ: INTEGER, num: n}
}

func exists(c *Client,v *Value, state *AppState) *Value {
	args := v.array[1:]
	var n int

	db.mu.RLock()
	for _, arg := range args {
		_, ok := db.store[arg.bulk]
		if ok {
			n++
		}
	}
	db.mu.RUnlock()

	return &Value{typ: INTEGER, num: n}
}

func keys(c *Client,v *Value, state *AppState) *Value {
	args := v.array[1:]
	if len(args) > 1 {
		return &Value{typ: ERROR, err: "ERR invalid number of arguments for 'KEYS' command"}
	}
	pattern := args[0].bulk

	db.mu.RLock()
	var matches []string
	for key := range db.store {
		matched, err := filepath.Match(pattern, key)
		if err != nil {
			log.Printf("error matching keys: (pattern: %s), (key: %s) - %v", pattern, key, err)
			continue
		}

		if matched {
			matches = append(matches, key)
		}
	}
	db.mu.RUnlock()

	reply := Value{typ: ARRAY}

	for _, m := range matches {
		reply.array = append(reply.array, Value{typ: BULK, bulk: m})
	}
	return &reply
}

func expire( c *Client,v *Value, state *AppState) *Value {
	args := v.array[1:]
	if len(args) != 2 {
		return &Value{typ: ERROR, err: "ERR invalid number of arguments for 'EXPIRE' command"}
	}

	k := args[0].bulk
	exp := args[1].bulk

	expSecs, err := strconv.Atoi(exp)
	if err != nil {
		return &Value{typ: ERROR, err: "ERR invalid expiry value"}
	}

	db.mu.RLock()
	key, ok := db.store[k]
	if !ok {
		return &Value{typ: INTEGER, num: 0}
	}
	key.Exp = time.Now().Add(time.Second * time.Duration(expSecs))
	db.mu.RUnlock()

	return &Value{typ: INTEGER, num: 1}
}

func ttl( c *Client,v *Value, state *AppState) *Value {
	args := v.array[1:]
	if len(args) != 1 {
		return &Value{typ: ERROR, err: "ERR invalid number of arguments for 'TTL' command"}
	}

	k := args[0].bulk

	db.mu.RLock()
	item, ok := db.store[k]
	if !ok {
		return &Value{typ: INTEGER, num: -2}
	}
	exp := item.Exp
	db.mu.RUnlock()
   
	if exp.Unix() == (time.Time{}).Unix(){
		return &Value{typ: INTEGER, num: -1}
	}

	expSecs:=int(time.Until(exp).Seconds())
	if expSecs<=0{
		db.mu.Lock()
		db.Delete(k)
		db.mu.Unlock()
		return &Value{typ: INTEGER, num: -2}
	}

	return &Value{typ: INTEGER, num: expSecs}

	
}

func multi(c *Client, v *Value, state *AppState) *Value {
	if state.tx != nil {
		return &Value{typ: ERROR, err: "ERR MULTI calls can not be nested"}
	}

	state.tx = NewTransaction()

	return &Value{typ: STRING, str: "OK"}
}

func _exec(c *Client, v *Value, state *AppState) *Value {
	if state.tx == nil {
		return &Value{typ: ERROR, err: "ERR EXEC without MULTI"}
	}

	replies := make([]Value, len(state.tx.cmds))
	for i, cmd := range state.tx.cmds {
		reply := cmd.handler(c, cmd.v, state)
		replies[i] = *reply
	}

	reply := Value{typ: ARRAY, array: replies}

	state.tx = nil

	return &reply
}

func discard(c *Client, v *Value, state *AppState) *Value {
	if state.tx == nil {
		return &Value{typ: ERROR, err: "ERR DISCARD without MULTI"}
	}

	state.tx = nil
	return &Value{typ: STRING, str: "OK"}
}

