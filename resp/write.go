package resp

import (
	"bufio"
	"fmt"
	"strconv"
)

/*
1. In this write.go file, we serialize the response accroding to RESP and write back to the client
*/
func WriteResp(w *bufio.Writer, v Value) error {
	bytes, err := v.Marshal()
	if err != nil {
		return err
	}

	_, err = w.Write(bytes)
	if err != nil {
		return err
	}
	return w.Flush()
}

func (v *Value) Marshal() ([]byte, error) {
	switch v.Typ {
	case "string":
		return v.marshalString()
	case "bulk":
		return v.marshalBulkString()
	case "array":
		return v.marshalArray()
	case "int":
		return v.marshalInt()
	case "null":
		return v.marshalNull()
	case "error":
		return v.marshalError()
	default:
		return []byte{}, fmt.Errorf("Unexpected type: %s", v.Typ)
	}
}

func (v *Value) marshalString() ([]byte, error) {
	// eg:  +OK\r\n
	return []byte(STRING + v.Str + "\r\n"), nil
}

func (v *Value) marshalError() ([]byte, error) {
	// eg: -Error message\r\n
	return []byte(ERROR + v.Err.Error() + "\r\n"), nil
}

func (v *Value) marshalNull() ([]byte, error) {
	// eg: $-1\r\n
	return []byte("$-1\r\n"), nil
}

func (v *Value) marshalBulkString() ([]byte, error) {
	// eg: $5\r\nhello\r\n
	// we have to use strconv.Itoa or fmt.Sprintf and not string() to convert int to string, as string() converts int to the ASCII equivalent character instead.
	return []byte(BULK + strconv.Itoa(len(v.Bulk)) + "\r\n" + v.Bulk + "\r\n"), nil
}

func (v *Value) marshalInt() ([]byte, error) {
	// eg:  :1000\r\n
	return []byte(INTEGER + strconv.Itoa(v.Integer) + "\r\n"), nil
}

func (v *Value) marshalArray() ([]byte, error) {
	// eg: *2\r\n$5\r\nhello\r\n$5\r\nworld\r\n
	result := []byte(ARRAY + strconv.Itoa(len(v.Array)) + "\r\n")

	// {array  0  [{bulk  0 set []} {bulk  0 hello []} {bulk  0 world []}]}
	for _, value := range v.Array {
		b, err := value.Marshal()
		if err != nil {
			return nil, err
		}

		result = append(result, b...)
	}
	return result, nil
}
