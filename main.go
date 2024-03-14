package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	fmt.Println("Listening on port 6379")
	// Start the server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Accept the connection
	conn, err := l.Accept()
	if err != nil {
		log.Fatalln(err.Error())
	}
	// Close the connection after we are done so that there's no resource leakage
	defer conn.Close()

	// Create a buffer (assign memory block) to store the contents from the connection
	buffer := make([]byte, 512)
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			log.Fatalln(err.Error())
		}
		// +OK\r\n is RESP notation, Redis client expects the string to be in this format
		conn.Write([]byte("+OK\r\n"))
	}
}
