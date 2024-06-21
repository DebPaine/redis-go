package main

import (
	"fmt"
	"log"
	"net"
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
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Keep reading from the connection perpetually
	for {
		buff := make([]byte, 1024)
		_, err := conn.Read(buff)
		if err != nil {
			log.Fatal(err.Error())
		}
		conn.Write([]byte("+OK\r\n"))
	}
}
