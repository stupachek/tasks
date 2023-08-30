package cli

import (
	"bytes"
	"epam/database"
	"epam/format"
	"epam/model"
	"fmt"
	"io"
	"os"
	"strconv"
)

var FileName string

func ReadCL() {
	var op string
	fmt.Print("> ")
	fmt.Scanf("%s\n", &op)
	switch op {
	case "INSERT":
		var name string
		var age int
		var active bool
		var mass float64
		var book string
		var fileName string
		books := make([]string, 0)
		fmt.Print("> File: ")
		fmt.Scanf("%s\n", &fileName)
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
		fmt.Printf("age %v", age)
		fmt.Printf("active %v", active)
		fmt.Print("> Mass: ")
		for _, err := fmt.Scanf("%f\n", &mass); err != nil; _, err = fmt.Scanf("%f\n", &mass) {
			fmt.Print("Mass must be floating-point number,! Try again: ")
		}
		fmt.Printf("mass %v", mass)
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
		}, fileName)
	case "DELETE":
		var index int
		for {
			var indexstr string
			var err error
			fmt.Print("> Index: ")
			fmt.Scanf("%s\n", &indexstr)
			index, err = strconv.Atoi(indexstr)
			if err != nil {
				fmt.Print("Index must be integer number! Try again:\n")
				continue
			}
			break
		}
		fmt.Println(index)
	case "CREATETABLE":
		fmt.Print("> Name: ")
		fmt.Scanf("%s\n", &FileName)
		_, err := os.Create(FileName)
		if err != nil {
			fmt.Printf("can't create file: %v", err)
		}
	case "SELECT":
		b := Reader()
		fmt.Printf("%b\n", b)
		usersModel := database.Users(b)
		fmt.Printf("%v\n", usersModel)
		format.PrepareUsers(usersModel)
		fmt.Println(format.FormatUsers((usersModel)))
	case "q":
		return
	default:
		fmt.Println("Unknown command!Try again.")
	}
	ReadCL()

}
func Reader() io.Reader {
	fmt.Print("> Name: ")
	fmt.Scanf("%s\n", &FileName)
	FileName = "1"
	file, err := os.OpenFile(FileName, os.O_RDONLY, 0777)
	if err != nil {
		fmt.Printf("can't open file: %v", err)
	}
	users := make([]byte, 500)
	n, err := file.Read(users)
	fmt.Printf("n %v, u %v", n, users)
	if err != nil {
		fmt.Printf("can't read file: %v", err)
	}
	return bytes.NewBuffer(
		users,
	)
}

func Test() {
	b := Reader()
	fmt.Printf("%b\n", b)
	usersModel := database.Users(b)
	fmt.Printf("%v\n", usersModel)
	format.PrepareUsers(usersModel)
	fmt.Println(format.FormatUsers((usersModel)))
}
