package resp

import (
	"bufio"
	"fmt"
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
	switch v.typ {
	case "string":
		return v.marshalString()
	case "bulk":
		return v.marshalBulkString()
	case "array":
	case "int":
	case "null":
	case "error":
	default:
		return []byte{}, fmt.Errorf("Unexpected type: %s", v.typ)
	}
	return []byte{}, fmt.Errorf("Unexpected type: %s", v.typ)
}

func (v *Value) marshalString() ([]byte, error) {
	// eg:  +OK\r\n
	return []byte(STRING + v.str + "\r\n"), nil
}

func (v *Value) marshalBulkString() ([]byte, error) {
	// eg: $5\r\nhello\r\n
	return []byte(BULK + string(len(v.bulk)) + "\r\n" + v.bulk + "\r\n"), nil
}
