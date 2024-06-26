package main

import (
	"fmt"
	"log"
	"net"

	"redis-go/resp"
)

/*
1. User enters the command in redis-cli [DONE]
2. We parse it using resp/read.go according to RESP, we store it in Value{} [DONE]
3. We interpret the command and take the necessary steps
4. We write back the response to the client using resp/write.go [IN-PROGRESS]

We are basically doing the following:
1. We first populate Value{} with the relevant command that the user has entered
2. We then use this Value{} to read the commands and execute it accordingly
3. After executing the command, we return a new Value{} based on the appropriate response

Eg:
Command:
SET hello world

RESP representation:
*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n

Value:
{array  0  [{bulk  0 set []} {bulk  0 hello []} {bulk  0 world []}]}

Return:
Value{typ: "string", str: "OK"}

RESP response:
+OK\r\n (simple string response)

CLI response:
OK

Basically, we initially convert RESP bytes to Value{} struct (reader), then convert Value{} struct back to RESP bytes (writer).
*/

func main() {
	// Listen on TCP port 6379
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer listener.Close()

	fmt.Println("Redis server started on port 6379")

	// Accept new connections from the Listener perpetually
	for {
		// This will block till a new connection is received
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err.Error())
		}
		// Handle the new connection logic using a new goroutine
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Keep reading from the connection perpetually
	for {
		rsp := resp.NewResp(conn)
		// value will be of Value struct, which will have values in it's array field
		value, err := rsp.Read()
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Println(value)
		conn.Write([]byte("+OK\r\n"))
	}
}
