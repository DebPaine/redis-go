package resp

import (
	"fmt"
	"strings"
)

func ExecuteCommand(v Value) (Value, error) {
	// eg: {array  0  [{bulk  0 set [] <nil>} {bulk  0 hello [] <nil>} {bulk  0 world [] <nil>}] <nil>}
	command, args := strings.ToLower(v.Array[0].Bulk), v.Array[1:]

	switch command {
	case "ping":
		return ping(args)
	case "set":
	case "get":
	default:
		return Value{}, fmt.Errorf("Unexpected command: %s", command)
	}
	return Value{}, fmt.Errorf("Unexpected command: %s", command)
}

func ping(args []Value) (Value, error) {
	if len(args) == 0 {
		return Value{Typ: "string", Str: "PONG"}, nil
	}
	return Value{Typ: "string", Str: args[0].Bulk}, nil
}
