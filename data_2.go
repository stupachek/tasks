package main

import (
	"encoding/base64"
	"fmt"
	"math"
	"strings"

	"golang.org/x/exp/slices"
)

type UserType int8

const (
	Organizer UserType = iota
	Critic
	CasualReader
	NonFictionBuff
	NewbieReader
)

func (u UserType) String() string {
	switch u {
	case Organizer:
		return "Organizer"
	case Critic:
		return "Critic"
	case CasualReader:
		return "Casual Reader"
	case NonFictionBuff:
		return "Non-Fiction Buff"
	case NewbieReader:
		return "Newbie Reader"
	default:
		return "Unknown"
	}
}

type User struct {
	Name     string
	Age      int
	Active   bool
	Mass     float64
	UserType UserType
	Books    []string
}

func Users() []User {
	return []User{
		{
			"John Doe",
			30,
			true,
			80.0,
			CasualReader,
			[]string{"Harry Potter", "1984"},
		},
		{
			"Jake Doe",
			20,
			false,
			60.0,
			NewbieReader,
			[]string{},
		},
		{
			" Jane Doe ",
			150,
			true,
			.75,
			NonFictionBuff,
			[]string{"Harry Potter", "Game of Thrones"},
		},
		{
			"\t",
			-10,
			true,
			8000.0,
			9,
			[]string{"Harry Potter"},
		},
		{
			"Vm0weE5GVXhUblJWV0dSUFZtMW9WVll3WkRSV1ZteDBaRVYwVmsxWGVGWlZiVEZIWVd4S2MxTnNiR0ZXVm5Cb1ZsVmFWMVpWTVVWaGVqQTk=\nVm0weE5GVXhUblJWV0dSUFZtMW9WVll3WkRSV1ZteDBaRVYwVmsxWGVGWlZiVEZIWVd4S2MxTnNiR0ZXVm5Cb1ZsVmFWMVpWTVVWaGVqQTk=",
			0,
			true,
			0,
			Organizer,
			[]string{"The Hunger Games"},
		},
		{
			"\x00\x10\x20\x30\x40\x50\x60\x70",
			0,
			true,
			0,
			Critic,
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

const nameLength = 8
const maxTypeLength = len("Non-Fiction Buff")

var (
	maxNameLength   = len("Name")
	maxBookLength   = len("Books")
	maxAgeLength    = len("Age")
	maxActiveLenght = len("Active")
	maxMassLength   = len("Mass")
)

func formatUsers(users []User) string {
	usersStr := make([]string, len(users))
	spacesBeforeBook := "\n" + strings.Join([]string{strings.Repeat(" ", maxNameLength), strings.Repeat(" ", maxAgeLength), strings.Repeat(" ", maxActiveLenght), strings.Repeat(" ", maxMassLength), strings.Repeat(" ", maxTypeLength)}, " | ") + " | "
	for i, user := range users {
		name := fmt.Sprintf("%*s", maxNameLength, users[i].Name)
		age := fmt.Sprintf("%+*d", maxAgeLength, user.Age)
		active := strings.Replace(fmt.Sprint(user.Active), "true", "+", 1)
		active = strings.Replace(active, "false", "-", 1)
		active = fmt.Sprintf("%*s", maxActiveLenght, active)
		mass := fmt.Sprintf("% *.2f", maxMassLength, user.Mass)
		userType := fmt.Sprintf("%*s", maxTypeLength, user.UserType)
		books := strings.Join(users[i].Books, spacesBeforeBook)
		usersStr[i] = strings.Join([]string{name, age, active, mass, userType, books}, " | ")
	}
	delimiter := strings.ReplaceAll(spacesBeforeBook, " ", "_") + strings.Repeat("_", maxBookLength) + "\n"
	header := fmt.Sprintf("%*s | %*s | %*s | %*s | %*s | %s", maxNameLength, "Name", maxAgeLength, "Age", maxActiveLenght, "Active", maxMassLength, "Mass", maxTypeLength, "Type", "Book") + delimiter
	return header + strings.Join(usersStr, delimiter)
}

func prepareUsers(users []User) {
	for i, user := range users {
		users[i].Name = strings.Trim(user.Name, " ")
		for len(users[i].Name) > nameLength {
			decoded, err := base64.StdEncoding.DecodeString(strings.Split(users[i].Name, "\n")[0])
			if err != nil {
				break
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
}

func averageAge(users []User) map[string]int {
	bookCounter := make(map[string]int)
	averageBookAge := make(map[string]int)
	for _, user := range users {
		for _, book := range user.Books {
			bookCounter[book]++
			averageBookAge[book] += user.Age
		}
	}
	for book := range averageBookAge {
		averageBookAge[book] /= bookCounter[book]
	}
	return averageBookAge

}

func strAvetageAge(averageBookAge map[string]int) string {
	averageTitle := "Avegare age"
	stringAverageAge := fmt.Sprintf("%*s | %s\n", maxBookLength, "Book", averageTitle)
	delimiter := strings.Repeat("_", maxBookLength) + " | " + strings.Repeat("_", len(averageTitle)) + "\n"
	stringAverageAge += delimiter
	for book := range averageBookAge {
		stringAverageAge += fmt.Sprintf("%*s | %*d\n", maxBookLength, book, len(averageTitle), averageBookAge[book])
		stringAverageAge += delimiter
	}
	return stringAverageAge
}

func nearestMass(users []User, target float64) User {
	slices.SortFunc(users, func(a User, b User) bool {
		return a.Mass < b.Mass
	})
	i, ok := slices.BinarySearchFunc(users, target, func(u User, mass float64) int {
		if u.Mass < mass {
			return -1
		}
		if u.Mass > mass {
			return 1
		}
		return 0
	})
	if !ok && i != 0 {
		if math.Abs(target-users[i-1].Mass) < math.Abs(target-users[i].Mass) {
			return users[i-1]
		}
	}
	return users[i]
}
func main() {
	users := Users()
	prepareUsers(users)
	fmt.Println(formatUsers((users)))
	averageAge := averageAge(users)
	fmt.Println(strAvetageAge(averageAge))
	slices.SortFunc(users, func(a User, b User) bool {
		var bookA, bookB int
		for _, book := range a.Books {
			bookA += averageAge[book]
		}
		for _, book := range b.Books {
			bookB += averageAge[book]
		}
		return bookA < bookB
	})
	fmt.Println("________Sort by sum of books average age_______")
	fmt.Println(formatUsers((users)))
	mass := 80.0
	fmt.Printf("________Nearest to mass %.1f_______\n", mass)
	fmt.Println(formatUsers([]User{nearestMass(users, mass)}))

}
