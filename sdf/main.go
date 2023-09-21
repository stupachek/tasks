package main

import (
	"fmt"
	"net"
)

func main() {
	conn, _ := net.Dial("tcp", "127.0.0.1:8080")
	b := make([]byte, 25)
	conn.Read(b)
	fmt.Println(string(b))
	conn.Write([]byte("INSERT"))
	conn.Close()
}
