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

// store and retrieve {key: value} pairs using a hashmap
var (
	setMap     = map[string]string{}
	setRWMutex = sync.RWMutex{}
)

// store and retrieve {key1: {key2: value}} pairs using a hashmap of hashmap
var (
	hsetMap     = map[string]map[string]string{}
	hsetRWMutex = sync.RWMutex{}
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
	case "hset":
		return hset(args)
	case "hget":
		return hget(args)
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

	setRWMutex.Lock()
	setMap[key] = value
	setRWMutex.Unlock()

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

	setRWMutex.RLock()
	value, ok := setMap[key]
	setRWMutex.RUnlock()
	if !ok {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: value}
}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{
			Typ: "error",
			Err: fmt.Errorf("Incorrect no. of arguments for HSET operation"),
		}
	}

	key1 := args[0].Bulk
	key2 := args[1].Bulk
	value := args[2].Bulk

	hsetRWMutex.Lock()
	if _, ok := hsetMap[key1]; !ok {
		hsetMap[key1] = map[string]string{}
	}
	hsetMap[key1][key2] = value
	hsetRWMutex.Unlock()

	return Value{Typ: "string", Str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{
			Typ: "error",
			Err: fmt.Errorf("Incorrect no. of arguments for HGET operation"),
		}
	}

	key1 := args[0].Bulk
	key2 := args[1].Bulk

	hsetRWMutex.RLock()
	value, ok := hsetMap[key1][key2]
	hsetRWMutex.RUnlock()
	if !ok {
		return Value{Typ: "null"}
	}

	return Value{Typ: "bulk", Bulk: value}
}
