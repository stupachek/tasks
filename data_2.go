package main

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type User struct {
	Name   string
	Age    int
	Active bool
	Mass   float64
	Books  []string
}

type UserString struct {
	User
	Row string
}

func Users() []User {
	return []User{
		{
			"John Doe",
			30,
			true,
			80.0,
			[]string{"Harry Potter", "1984"},
		},
		{
			"Jake Doe",
			20,
			false,
			60.0,
			[]string{},
		},
		{
			" Jane Doe ",
			150,
			true,
			.75,
			[]string{"Harry Potter", "Game of Thrones"},
		},
		{
			"\t",
			-10,
			true,
			8000.0,
			[]string{"Harry Potter"},
		},
		{
			"Vm0weE5GVXhUblJWV0dSUFZtMW9WVll3WkRSV1ZteDBaRVYwVmsxWGVGWlZiVEZIWVd4S2MxTnNiR0ZXVm5Cb1ZsVmFWMVpWTVVWaGVqQTk=\nVm0weE5GVXhUblJWV0dSUFZtMW9WVll3WkRSV1ZteDBaRVYwVmsxWGVGWlZiVEZIWVd4S2MxTnNiR0ZXVm5Cb1ZsVmFWMVpWTVVWaGVqQTk=",
			0,
			true,
			0,
			[]string{"The Hunger Games"},
		},
		{
			"\x00\x10\x20\x30\x40\x50\x60\x70",
			0,
			true,
			0,
			[]string{"Moby Dick", "It", "The Green Mile"},
		},
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

const nameLength = 10

func main() {

	users := Users()
	usersStr := make([]string, len(users))
	maxNameLength := len("Name")
	maxBookLength := len("Books")
	maxAgeLength := len("Age")
	maxActiveLenght := len("Active")
	maxMassLength := len("Mass")
	for i, user := range users {
		users[i].Name = strings.Trim(user.Name, " ")
		for len(users[i].Name) > nameLength {
			decoded, err := base64.StdEncoding.DecodeString(strings.Split(users[i].Name, "\n")[0])
			if err != nil {
				panic(err)
			}
			users[i].Name = string(decoded)
		}
		users[i].Name = fmt.Sprintf("%q", users[i].Name)
		maxNameLength = max(maxNameLength, len(users[i].Name))
		maxAgeLength = max(maxAgeLength, len(fmt.Sprintf("%+d", user.Age)))
		maxMassLength = max(maxMassLength, len(fmt.Sprintf("%f", user.Mass)))
		for j, book := range users[i].Books {
			users[i].Books[j] = fmt.Sprintf("%q", strings.TrimSpace(book))
			maxBookLength = max(maxBookLength, len(users[i].Books[j]))
		}

	}
	spacesBeforeBook := "\n" + strings.Join([]string{strings.Repeat(" ", maxNameLength), strings.Repeat(" ", maxAgeLength), strings.Repeat(" ", maxActiveLenght), strings.Repeat(" ", maxMassLength-1)}, " | ") + "| "
	for i, user := range users {
		name := fmt.Sprintf("%*s", maxNameLength, users[i].Name)
		age := fmt.Sprintf("%+*d", maxAgeLength, user.Age)
		active := strings.Replace(fmt.Sprint(user.Active), "true", "+", 1)
		active = strings.Replace(active, "false", "-", 1)
		active = fmt.Sprintf("%*s", maxActiveLenght, active)
		mass := fmt.Sprintf("% *.2f", maxMassLength-2, user.Mass)
		for j, book := range users[i].Books {
			users[i].Books[j] = fmt.Sprintf("%*s", maxBookLength, book)
		}
		books := strings.Join(users[i].Books, spacesBeforeBook)
		usersStr[i] = strings.Join([]string{name, age, active, mass, books}, " | ")
	}
	fmt.Printf("%*s | %*s | %*s | %*s | %*s\n", maxNameLength, "Name", maxAgeLength, "Age", maxActiveLenght, "Active", maxMassLength-2, "Mass", maxBookLength, "Book")
	fmt.Println(strings.Join(usersStr, "\n"))

}
