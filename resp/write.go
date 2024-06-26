package resp

/*
1. In this write.go file, we serialize the response accroding to RESP and write back to the client
*/

func (v *Value) Marshal() []byte {
	return []byte("OK")
}
