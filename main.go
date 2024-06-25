package main

import (
	"fmt"
	"log"
	"net"
	"redis-go/resp"
)

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
