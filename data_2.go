package main

import (
	"epam/cli"
	"time"
)

func main() {
	n := cli.NewConnection()
	cli.ReadCL(&n)
	time.Sleep(time.Millisecond * 100)
}
