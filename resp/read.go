package resp

import (
	"bufio"
	"fmt"
	"io"
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
      inputType: "array",
      str: "",
      num: 0,
      bulk: "",
      array: [
              Value{inputType: "bulk", str: "", num: 0, bulk: "hello", array: []},
              Value{inputType: "bulk", str: "", num: 0, bulk: "world", array: []}
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
	inputType string  // type of input, eg: array, bulkstring, simplestring, integer
	str       string  // simple strings
	num       int     // integer values
	bulk      string  // bulk strings
	array     []Value // array values
}

type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	// bufio.NewReader returns a pointer, hence "reader" field is of type *bufio.Reader
	return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) Read() (Value, error) {
	// eg: *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
	for {
		// line is the string till \r\n, which is *2 in the first iteration
		line, err := r.reader.ReadString('\n')
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
		case ARRAY:
			v, err := r.readArray(line)
			if err != nil {
				return Value{}, err
			}
			return v, err
		default:
			return Value{}, fmt.Errorf("Unexpected type: %s", _type)
		}
	}
	return Value{}, nil
}

func (r *Resp) readArray(line string) (Value, error) {
	v := Value{}
	v.inputType = "array"

	arrayLength, err := strconv.Atoi(line[1:])
	if err != nil {
		return Value{}, err
	}

	// we are iterating till arrayLength since we have that many elements in the array
	for i := 0; i < arrayLength; i++ {
		value, err := r.readBulkString()
		if err != nil {
			return Value{}, err
		}

		v.array = append(v.array, value)
	}
	return v, nil
}

func (r *Resp) readBulkString() (Value, error) {
	// eg: $5\r\nhello\r\n
	v := Value{}

	// Read the bytes till \n
	line, err := r.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	// Remove empty spaces like \r\n
	line = strings.TrimSpace(line)
	// Get the first byte which is the type
	_type := string(line[0])
	if _type != BULK {
		return Value{}, fmt.Errorf("Expected %v, got %v", BULK, _type)
	}
	v.inputType = "bulk"

	// Convert the length of the string to int
	lineLength, err := strconv.Atoi(line[1:])
	if err != nil {
		return Value{}, err
	}

	// Get the actual data in the bulk string
	bulkString := make([]byte, lineLength)
	_, err = r.reader.Read(bulkString)
	if err != nil {
		return Value{}, err
	}
	v.bulk = string(bulkString)

	// eg: $5\r\nhello\r\n, after reading "hello" we have to read \r\n too so that we can move to the next bulkstring
	_, err = r.reader.ReadString('\n')
	if err != nil {
		return Value{}, err
	}

	return v, nil
}