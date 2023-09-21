package server

import (
	"fmt"
	"net"
	"time"
)

func Listen() {
	c := make(chan int)
	go func() {
		defer func() {
			c <- 0
		}()
		l, _ := net.Listen("tcp", ":8080")
		for {
			conn, _ := l.Accept()
			// conn.Read()
			conn.Write([]byte(`Write command!`))
			b := make([]byte, 25)
			conn.Read(b)
			fmt.Print(string(b))
			conn.Close()
		}
		//
	}()
	time.Sleep(time.Millisecond)
	go func() {
		defer func() {
			c <- 0
		}()
		conn, _ := net.Dial("tcp", "127.0.0.1:8080")
		b := make([]byte, 25)
		conn.Read(b)
		conn.Write([]byte(`qwerty`))
		fmt.Println(string(b))
		conn.Close()
	}()
	<-c
	<-c
}
