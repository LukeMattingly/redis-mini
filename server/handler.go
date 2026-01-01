package main

import "sync"

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
}

var SETs = map[string]string{}
var SETsMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{Type: TypeError, Str: "wrong number of args"}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	SETsMu.Lock()
	SETs[key] = value
	SETsMu.Unlock()

	return Value{Type: TypeSimpleString, Str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{Type: TypeSimpleString, Str: "wrong number of args"}
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
