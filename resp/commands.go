package resp

var CommandHandler = map[string]func([]Value) Value{
	"PING": ping,
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "string", Str: "PONG"}
	}
	return Value{Typ: "string", Str: args[0].Bulk}
}
