package main

import "sync"

const (
	OK           = "OK"
	WRONGNUMARGS = "Wrong number of args"
)

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Type: TypeSimpleString, Str: "PONG"}
	}
	return Value{Type: TypeSimpleString, Str: args[0].Bulk}
}

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
	"SET":  set,
	"GET":  get,
	"HSET": hset,
	"HGET": hget,
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{Type: TypeError, Str: WRONGNUMARGS}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{Type: TypeSimpleString, Str: OK}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{Type: TypeSimpleString, Str: WRONGNUMARGS}
	}

	key := args[0].Bulk

	SETsMu.RLock()
	value, ok := SETs[key]
	SETsMu.RUnlock()

	if !ok {
		return Value{Type: TypeNull, Str: "null"}
	}

	return Value{Type: TypeBulkString, Bulk: value}
}

var HSETs = map[string]map[string]string{}
var HSETsMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{Type: TypeError, Str: WRONGNUMARGS}
	}

	hash := args[0].Bulk
	key := args[1].Bulk
	value := args[2].Bulk

	HSETsMu.Lock()
	if _, ok := HSETs[hash]; !ok {
		HSETs[hash] = map[string]string{}
	}
	HSETs[hash][key] = value
	HSETsMu.Unlock()

	return Value{Type: TypeSimpleString, Str: OK}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{Type: TypeNull, Str: WRONGNUMARGS}
	}

	hash := args[0].Bulk
	key := args[1].Bulk

	HSETsMu.RLock()
	value, ok := HSETs[hash][key]
	HSETsMu.RUnlock()

	if !ok {
		return Value{Type: TypeNull}
	}

	return Value{Type: TypeBulkString, Bulk: value}
}
