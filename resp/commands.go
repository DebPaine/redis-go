package resp

import (
	"fmt"
	"strings"
	"sync"
)

/*
Use sync.Mutex when you need simple mutual exclusion for critical sections.
Use sync.RWMutex when you have more reads than writes and want to allow concurrent read access while still ensuring exclusive write access.

Concurrency:
RLock allows multiple concurrent readers.
Lock allows only one writer and no concurrent readers or writers.

Blocking Behavior:
RLock will block if a write lock is currently held.
Lock will block if any read or write lock is currently held.

Explanation:
If an RLock (read lock) is held by any goroutine, a Lock (write lock) cannot be acquired by another goroutine. This is one of the fundamental properties of sync.RWMutex in Go: read locks and write locks are mutually exclusive to ensure data consistency.
*/

var (
	setMap  = map[string]string{}
	rwMutex = sync.RWMutex{}
)

func ExecuteCommand(v Value) Value {
	// eg: {array  0  [{bulk  0 set [] <nil>} {bulk  0 hello [] <nil>} {bulk  0 world [] <nil>}] <nil>}
	command, args := strings.ToLower(v.Array[0].Bulk), v.Array[1:]

	switch command {
	case "ping":
		return ping(args)
	case "set":
		return set(args)
	case "get":
		return get(args)
	default:
		return Value{Typ: "null"}
	}
}

func ping(args []Value) Value {
	if len(args) == 0 {
		return Value{Typ: "string", Str: "PONG"}
	}
	return Value{Typ: "string", Str: args[0].Bulk}
}

func set(args []Value) Value {
	// eg: {array  0  [{bulk  0 set [] <nil>} {bulk  0 hello [] <nil>} {bulk  0 world [] <nil>}] <nil>}
	if len(args) != 2 {
		return Value{
			Typ: "error",
			Err: fmt.Errorf("Incorrect no. of arguments for SET operation"),
		}
	}

	key := args[0].Bulk
	value := args[1].Bulk

	rwMutex.Lock()
	setMap[key] = value
	rwMutex.Unlock()

	return Value{Typ: "string", Str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{
			Typ: "error",
			Err: fmt.Errorf("Incorrect no. of arguments for GET operation"),
		}
	}

	key := args[0].Bulk

	rwMutex.RLock()
	value, ok := setMap[key]
	if !ok {
		return Value{Typ: "null"}
	}
	rwMutex.RUnlock()

	return Value{Typ: "bulk", Bulk: value}
}
