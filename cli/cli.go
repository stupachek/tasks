package cli

import (
	"epam/database"
	"epam/format"
	"epam/model"
	"fmt"
	"net"
	"os"
	"strconv"
)

type UserInterface interface {
	Command() string
	File()
}

type CommandLine struct{}

func (c *CommandLine) Command() string {
	var op string
	fmt.Print("> ")
	fmt.Scanf("%s\n", &op)
	return op
}

func (c *CommandLine) File() {
	fmt.Print("> File: ")
	fmt.Scanf("%s\n", &FileName)
}

type Net struct {
	conn net.Conn
}

func NewConnection() Net {
	l, _ := net.Listen("tcp", ":8080")
	for {
		conn, _ := l.Accept()
		n := Net{
			conn: conn,
		}
		return n
	}
}

func (n *Net) Command() string {
	n.conn.Write([]byte("> "))
	op := make([]byte, 10)
	num, _ := n.conn.Read(op)
	return string(op[:num])
}

func (n *Net) File() {
	n.conn.Write([]byte("> File: "))
	file := make([]byte, 10)
	num, _ := n.conn.Read(file)
	FileName = string(file[:num])
}

var FileName string

func ReadCL(ui UserInterface) {
	op := ui.Command()
	switch op {
	case "INSERT":
		var name string
		var age int
		var active bool
		var mass float64
		var book string
		books := make([]string, 0)
		if FileName == "" {
			ui.File()
		}
		fmt.Print("> Name: ")
		fmt.Scanf("%s\n", &name)
		for {
			var agestr string
			var err error
			fmt.Print("> Age: ")
			fmt.Scanf("%s\n", &agestr)
			age, err = strconv.Atoi(agestr)
			if err != nil {
				fmt.Print("Age must be integer number! Try again:\n")
				continue
			}
			break
		}
		for {
			var activestr string
			fmt.Print("> Active: ")
			fmt.Scanf("%s\n", &activestr)
			switch activestr {
			case "true":
				active = true
			case "false":
				active = false
			default:
				fmt.Print("Age must be integer number! Try again:\n")
				continue
			}
			break
		}
		fmt.Print("> Mass: ")
		for _, err := fmt.Scanf("%f\n", &mass); err != nil; _, err = fmt.Scanf("%f\n", &mass) {
			fmt.Print("Mass must be floating-point number,! Try again: ")
		}
		fmt.Print("> Book: ")
		for fmt.Scanf("%s\n", &book); book != "q"; fmt.Scanf("%s\n", &book) {
			books = append(books, book)
			fmt.Print("> Book: ")
		}
		database.Insert(model.User{
			Name:   name,
			Age:    uint64(age),
			Active: active,
			Mass:   mass,
			Books:  books,
		}, FileName)
	case "CREATETABLE":
		fmt.Print("> Name: ")
		fmt.Scanf("%s\n", &FileName)
		_, err := os.Create(FileName)
		if err != nil {
			fmt.Printf("can't create file: %v", err)
		}
	case "SELECT":
		if FileName == "" {
			fmt.Print("> File name: ")
			fmt.Scanf("%s\n", &FileName)
		}
		file, err := os.OpenFile(FileName, os.O_RDONLY, 0777)
		if err != nil {
			fmt.Printf("can't open file: %v", err)
		}
		usersModel := database.Users(file)
		format.PrepareUsers(usersModel)
		fmt.Println(format.FormatUsers((usersModel)))
	case "q":
		return
	default:
		fmt.Println("Unknown command!Try again.")
	}
	ReadCL(ui)

}
