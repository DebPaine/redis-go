package resp

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	// "strconv"
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

// func (r *Resp) Read() (Value, error) {
// 	// eg: $5\r\nhello\r\n
// 	_type, err := r.reader.ReadByte()
// 	if err != nil {
// 		log.Fatalln(err.Error())
// 	}
//
// 	switch _type {
// 	case ARRAY:
//
// 	case BULK:
// 	default:
// 		fmt.Printf("Unknown type: %v", string(_type))
// 		return Value{}, nil
// 	}
// }

func (r *Resp) Read() (Value, error) {
	// eg: *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
	v := Value{}
	v.inputType = "array"

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
			arrayLength, err := strconv.Atoi(line[1:])
			if err != nil {
				return Value{}, err
			}

			for i := 0; i < arrayLength; i++ {
				value, err := r.readBulkString()
				if err != nil {
					return Value{}, err
				}

				v.array = append(v.array, value)
			}
		}
	}
	return v, nil
}

func (r *Resp) readBulkString() (Value, error) {
	return Value{}, nil
}

// func (r *Resp) readLine() (line []byte, n int, err error){
//   // Keep reading byte by byte till we reach \r\n and stop before that
//   // eg: $5\r\nhello\r\n, here we will only read till $5 then return
// 	for {
//     b, err := r.reader.ReadByte()
//     if err != nil {
//       return nil, 0, err  // nil is a valid return value for a slice
//     }
//     if b == '\r' {
//       break
//     }
//     n += 1
//     line = append(line, b)
//   }
//   return line, n, nil
// }
//
// func (r *Resp) readInteger() (x int, n int, err error){
//   line, n, err := r.readLine()
// 	if err != nil {
//     return 0, 0, err
//   }
//   i64, err := strconv.ParseInt(string(line), 10, 64)
//   if err != nil {
//     return 0, n, err
//   }
//   return int(i64), n, nil
// }
