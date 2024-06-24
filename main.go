package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"redis-go/resp"
	"strings"
)

func main() {
	input := "$3\r\nDeb\r\n"
	reader := bufio.NewReader(strings.NewReader(input))
	b, err := reader.ReadString('\n')
	fmt.Println(b[:len(b)-2])
	b, err = reader.ReadString('\n')
	fmt.Println(b[:len(b)-2])
	// b, err := reader.ReadByte()
	// fmt.Println("Type", string(b))
	//
	// size, _ := reader.ReadByte()
	// strSize, _ := strconv.ParseInt(string(size), 10, 64)
	// fmt.Println("Size", strSize)
	//
	// reader.ReadByte()
	// reader.ReadByte()
	//
	// buff := make([]byte, size)
	// reader.Read(buff)
	// fmt.Println("Text", string(buff))

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
		value, err := rsp.Read()
		if err != nil {
			log.Fatalln(err.Error())
		}

		fmt.Println(value)
		conn.Write([]byte("+OK\r\n"))
		// buff := make([]byte, 1024)
		// _, err := conn.Read(buff)
		// if err != nil {
		// 	// EOF error occurs when we close the connection, we break from the current infinite loop using break
		// 	if err == io.EOF {
		// 		fmt.Println("EOF error")
		// 		break
		// 	}
		// 	log.Fatal(err.Error())
		// }
		// conn.Write([]byte("+OK\r\n"))
	}
}
