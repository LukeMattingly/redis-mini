package main

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Type: TypeSimpleString, Str: "PONG"}
	}
	return Value{Type: TypeSimpleString, Str: args[0].Bulk}
}

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
}
