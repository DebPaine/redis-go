package resp

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

/*
In this read.go file, we deserialize the user input and parse the commands which follows RESP standards. Here, we are only accepting arrays and bulkstrings from user input.
Here is how it works:
1. We use the Read method to parse and see if it's an array type or not.
2. If it's an array type, we then parse it further using readArray method.
3. We further parse it in readArray method and then check if we have bulkstrings or not.
4. If yes, we then further parse the bulkstrings using readBulkString method.
5. We save the result of all this in Value struct.

Initially, the Value struct will be empty. As we parse the input further, we keep adding updating the necessary fields in the Value struct according to the user input. If we have
an array input, then we also update the array field in Value, which itself is of Value type. Here is how Value struct will look like once a user input is parsed:

User input:
set hello world

RESP representation of input:
*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n

Parsed input:
{array  0  [{bulk  0 set []} {bulk  0 hello []} {bulk  0 world []}]}

Parsed input with Value struct fields with values:
Value{
      typ: "array",
      str: "",
      num: 0,
      bulk: "",
      array: [
              Value{typ: "bulk", str: "", num: 0, bulk: "hello", array: []},
              Value{typ: "bulk", str: "", num: 0, bulk: "world", array: []}
            ]
    }

As you can see, the parsed input is Value struct with updated array field. The array field has multiple Value structs each of "bulk" string type.
*/

const (
	STRING  = "+"
	ERROR   = "-"
	INTEGER = ":"
	BULK    = "$"
	ARRAY   = "*"
)

// Value struct will hold the command entered by the user
type Value struct {
	typ     string  // type of input, eg: array, bulkstring, simplestring, integer
	str     string  // simple strings
	integer int     // integer values
	bulk    string  // bulk strings
	array   []Value // array values
	err     error   // error values
}

func ReadResp(r *bufio.Reader) (Value, error) {
	// eg: *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n

	for {
		// line is the string till \r\n, which is *2 in the first iteration
		line, err := r.ReadString('\n')
		if err != nil {
			return Value{}, err
		}

		// trim \r\n from the line
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			break
		}

		_type := string(line[0])

		switch _type {
		case STRING:
			return readString(line)
		case INTEGER:
			return readInteger(line)
		case ARRAY:
			return readArray(r, line)
		case BULK:
			return readBulkString(
				r,
				line,
			) // we need to send the reader in the function as just sending the line won't be enough since there are multiple \r\n in a bulkstring
		case ERROR:
			return readError(line)
		default:
			return Value{}, fmt.Errorf("Unexpected type: %s", _type)
		}
	}
	return Value{}, nil
}

func readString(s string) (Value, error) {
	// eg:  +OK\r\n
	return Value{typ: "string", str: s[1:]}, nil
}

func readInteger(s string) (Value, error) {
	// eg:  :1000\r\n
	integer, err := strconv.Atoi(s[1:])
	if err != nil {
		fmt.Println("Integer conversion error: ", err)
		return Value{}, err
	}

	return Value{typ: "int", integer: integer}, nil
}

func readError(line string) (Value, error) {
	// eg: -Error message\r\n
	return Value{typ: "error", err: fmt.Errorf(line[1:])}, nil
}

func readArray(r *bufio.Reader, line string) (Value, error) {
	// eg: *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n

	v := Value{typ: "array"}

	length, err := strconv.Atoi(line[1:])
	if err != nil {
		fmt.Println("Array length error: ", err)
		return Value{}, err
	}
	// we are iterating till arrayLength since we have that many elements in the array
	for i := 0; i < length; i++ {
		value, err := ReadResp(r)
		if err != nil {
			return Value{}, err
		}

		v.array = append(v.array, value)
	}
	return v, nil
}

func readBulkString(r *bufio.Reader, line string) (Value, error) {
	// eg: $5\r\nhello\r\n
	v := Value{typ: "bulk"}

	// Convert the length of the string to int
	length, err := strconv.Atoi(line[1:])
	if err != nil {
		fmt.Println("Bulkstring length error: ", err)
		return Value{}, err
	}

	// Get the actual data in the bulk string
	bulkString := make([]byte, length)
	_, err = r.Read(bulkString)
	if err != nil {
		return Value{}, err
	}

	v.bulk = string(bulkString)

	// eg: $5\r\nhello\r\n, after reading "hello" we have to read \r\n too so that we can move to the next bulkstring
	_, err = r.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	return v, nil
}
