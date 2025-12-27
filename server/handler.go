package main

func ping(args []Value) Value {
	return Value{Type: TypeSimpleString, Str: "PONG"}
}

var Handlers = map[string]func([]Value) Value{
	"PING": ping,
}
