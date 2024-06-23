package resp

import (
	"bufio"
	"io"
	"strconv"
)

const (
  STRING = "+"
  ERROR = "-"
  INTEGER = ":"
  BULK = "$"
  ARRAY = "*"
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
}

func NewResp(rd io.Reader) *Resp {
  // bufio.NewReader returns a pointer, hence "reader" field is of type *bufio.Reader 
  return &Resp{reader: bufio.NewReader(rd)}
}

func (r *Resp) readLine() (line []byte, n int, err error){
  // Keep reading byte by byte till we reach \r\n and stop before that
  // eg: $5\r\nhello\r\n, here we will only read till $5 then return
  for {
    b, err := r.reader.ReadByte()
    if err != nil {
      return nil, 0, err  // nil is a valid return value for a slice
    }
    if b == '\r' {
      break
    }
    n += 1
    line = append(line, b)
  }
  return line, n, nil
}

func (r *Resp) readInteger() (x int, n int, err error){
  line, n, err := r.readLine()
  if err != nil {
    return 0, 0, err
  }
  i64, err := strconv.ParseInt(string(line), 10, 64)
  if err != nil {
    return 0, n, err
  }
  return int(i64), n, nil
}
