package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

const (
	STRING  = "+"
	ERROR   = "-"
	INTEGER = ":"
	BULK    = "$"
	ARRAY   = "*"
)

// Value struct will hold the command entered by the user
type Value struct {
	inputType string
	str       string
	num       int
	bulk      string
	array     []Value
}

type Resp struct {
	reader *bufio.Reader
	writer *bufio.Writer
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
	v.num = lineLength

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
