package resp

var handlers = map[string]func([]Value) Value{
	"PING": ping,
}

func ping(_ []Value) Value {
	return Value{inputType: "string", str: "PONG"}
}
