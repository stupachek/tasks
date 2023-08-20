package main

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"strings"

	"golang.org/x/exp/slices"
)

const (
	maxUsers int = 8
)

var active uint8 = 0

type User struct {
	Name  string   // uint8(length) + [length]byte
	Age   uint64   // 1 bit bool (active field) + 63 bit uint (age field)
	Mass  float64  // regular float64
	Books []string // uint8(all books length) + [length]byte, all books come as a single comma-separated string
}

type UserFromBinary struct {
	Name      []byte
	ActiveAge uint64
	Mass      float64
}

func Users(r io.Reader) ([]User, error) {
	out := []User{}
	for i := 0; i < maxUsers; i++ {
		user, err := ReadUser(&r, i)
		if err != nil {
			break
		}
		out = append(out, user)
	}
	return out, nil
}

func ReadUser(r *io.Reader, position int) (User, error) {
	var nameLength uint8
	if err := binary.Read(*r, binary.BigEndian, &nameLength); err != nil {
		return User{}, err
	}
	name := make([]byte, nameLength)
	if err := binary.Read(*r, binary.BigEndian, &name); err != nil {
		return User{}, err
	}
	var ActiveAge uint64
	if err := binary.Read(*r, binary.BigEndian, &ActiveAge); err != nil {
		return User{}, err
	}
	bitmask := uint64(math.Exp2(64))
	activeUser := uint8(ActiveAge & bitmask >> (64 - 1))
	active ^= activeUser << (position)
	age := ActiveAge << 1
	age = age >> 1
	var mass float64
	if err := binary.Read(*r, binary.BigEndian, &mass); err != nil {
		return User{}, err
	}
	var bookLength uint8
	if err := binary.Read(*r, binary.BigEndian, &bookLength); err != nil {
		return User{}, err
	}
	book := make([]byte, bookLength)
	if err := binary.Read(*r, binary.BigEndian, &book); err != nil {
		return User{}, err
	}
	books := strings.Split(string(book), ",")
	return User{
		Name:  string(name),
		Age:   age,
		Mass:  mass,
		Books: books,
	}, nil
}

func Reader() io.Reader {
	return bytes.NewBuffer(
		[]byte{
			0x8, 0x4a, 0x6f, 0x68, 0x6e, 0x20, 0x44, 0x6f, 0x65, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1e, 0x40, 0x54, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x11, 0x48, 0x61, 0x72, 0x72, 0x79, 0x20, 0x50, 0x6f, 0x74, 0x74, 0x65, 0x72, 0x2c, 0x31, 0x39, 0x38, 0x34, 0x8, 0x4a, 0x61, 0x6b, 0x65, 0x20, 0x44, 0x6f, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0x40, 0x4e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x20, 0x4a, 0x61, 0x6e, 0x65, 0x20, 0x44, 0x6f, 0x65, 0x20, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x96, 0x3f, 0xe8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1c, 0x48, 0x61, 0x72, 0x72, 0x79, 0x20, 0x50, 0x6f, 0x74, 0x74, 0x65, 0x72, 0x2c, 0x47, 0x61, 0x6d, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x54, 0x68, 0x72, 0x6f, 0x6e, 0x65, 0x73, 0x1, 0x9, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5a, 0x40, 0xbf, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc, 0x48, 0x61,
			0x72, 0x72, 0x79, 0x20, 0x50, 0x6f, 0x74, 0x74, 0x65, 0x72, 0xd9, 0x56, 0x6d, 0x30, 0x77, 0x65, 0x45, 0x35, 0x47, 0x56, 0x58, 0x68, 0x55, 0x62, 0x6c, 0x4a, 0x57, 0x56, 0x30, 0x64, 0x53, 0x55, 0x46, 0x5a, 0x74, 0x4d, 0x57, 0x39, 0x57, 0x56, 0x6c, 0x6c, 0x33, 0x57, 0x6b, 0x52,
			0x53, 0x56, 0x31, 0x5a, 0x74, 0x65, 0x44, 0x42, 0x61, 0x52, 0x56, 0x59, 0x77, 0x56, 0x6d, 0x73, 0x78, 0x57, 0x47, 0x56, 0x47, 0x57, 0x6c, 0x5a, 0x69, 0x56, 0x45, 0x5a, 0x49, 0x57, 0x56, 0x64, 0x34, 0x53, 0x32, 0x4d, 0x78, 0x54, 0x6e, 0x4e, 0x69, 0x52, 0x30, 0x5a, 0x58, 0x56,
			0x6d, 0x35, 0x43, 0x62, 0x31, 0x5a, 0x73, 0x56, 0x6d, 0x46, 0x57, 0x4d, 0x56, 0x70, 0x57, 0x54, 0x56, 0x56, 0x57, 0x61, 0x47, 0x56, 0x71, 0x51, 0x54, 0x6b, 0x3d, 0xa, 0x56, 0x6d, 0x30, 0x77, 0x65, 0x45, 0x35, 0x47, 0x56, 0x58, 0x68, 0x55, 0x62, 0x6c, 0x4a, 0x57, 0x56, 0x30, 0x64, 0x53, 0x55, 0x46, 0x5a, 0x74, 0x4d, 0x57, 0x39, 0x57, 0x56, 0x6c, 0x6c, 0x33, 0x57, 0x6b, 0x52, 0x53, 0x56, 0x31, 0x5a, 0x74, 0x65, 0x44, 0x42, 0x61, 0x52, 0x56, 0x59, 0x77, 0x56, 0x6d, 0x73, 0x78, 0x57, 0x47, 0x56, 0x47, 0x57, 0x6c, 0x5a, 0x69, 0x56, 0x45, 0x5a, 0x49, 0x57, 0x56, 0x64, 0x34, 0x53, 0x32, 0x4d, 0x78, 0x54, 0x6e, 0x4e, 0x69, 0x52, 0x30, 0x5a, 0x58, 0x56, 0x6d, 0x35, 0x43, 0x62, 0x31, 0x5a, 0x73, 0x56, 0x6d, 0x46, 0x57, 0x4d, 0x56, 0x70, 0x57, 0x54, 0x56, 0x56, 0x57, 0x61, 0x47, 0x56, 0x71, 0x51, 0x54, 0x6b, 0x3d, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x54, 0x68, 0x65, 0x20, 0x48, 0x75, 0x6e, 0x67, 0x65, 0x72, 0x20, 0x47, 0x61, 0x6d, 0x65, 0x73, 0x8, 0x0, 0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b, 0x4d, 0x6f, 0x62, 0x79, 0x20, 0x44, 0x69, 0x63, 0x6b, 0x2c, 0x49, 0x74, 0x2c, 0x54, 0x68, 0x65, 0x20, 0x47, 0x72, 0x65, 0x65, 0x6e, 0x20, 0x4d, 0x69, 0x6c, 0x65,
		},
	)
}

const (
	nameLength = 8
	Name       = "Name"
	Books      = "Books"
	Age        = "Age"
	Active     = "Active (total %d)"
	Mass       = "Mass"
)

var (
	maxNameLength   = len(Name)
	maxBookLength   = len(Books)
	maxAgeLength    = len(Age)
	maxActiveLenght = len(Active)
	maxMassLength   = len(Mass)
)

func count(n uint8) int {
	counter := 0
	for i := 0; i < 8; i++ {
		counter += int(n & 1)
		n >>= 1
	}
	return counter
}

func delimiterPattern(vertical, horizontal string) string {
	return strings.Join([]string{strings.Repeat(vertical, maxNameLength), strings.Repeat(vertical, maxAgeLength), strings.Repeat(vertical, maxActiveLenght), strings.Repeat(vertical, maxMassLength-2), strings.Repeat(vertical, maxBookLength)}, horizontal)
}

func formatUsers(users []User) string {
	usersStr := make([]string, len(users))
	upperBound := "┏" + delimiterPattern("━", "━┳━") + "┓\n"
	spacesBeforeBook := delimiterPattern(" ", " ┃ ")
	spacesBeforeBook = spacesBeforeBook[:len(spacesBeforeBook)-maxBookLength]
	lowerBound := "\n┗" + delimiterPattern("━", "━┻━") + "┛"
	for i, user := range users {
		name := fmt.Sprintf("%*s", maxNameLength, users[i].Name)
		age := fmt.Sprintf("%*d", maxAgeLength, user.Age)
		activeMask := active & (1 << uint8(i))
		activestr := fmt.Sprintf("%*s", maxActiveLenght, "-")
		if activeMask > 0 {
			activestr = fmt.Sprintf("%*s", maxActiveLenght, "+")
		}
		mass := fmt.Sprintf("% *.2f", maxMassLength-2, user.Mass)
		bookFormat(user.Books)
		books := strings.Join(users[i].Books, "┃\n┃"+spacesBeforeBook)
		usersStr[i] = "┃" + strings.Join([]string{name, age, activestr, mass, books}, " ┃ ") + "┃"
	}
	delimiter := "\n┣" + delimiterPattern("━", "━╋━") + "┫\n"
	activeHeader := fmt.Sprintf(Active, count(active))
	header := fmt.Sprintf("┃%*s ┃ %*s ┃ %*s ┃ %*s ┃ %*s┃", maxNameLength, Name, maxAgeLength, Age, maxActiveLenght, activeHeader, maxMassLength-2, Mass, maxBookLength, Books) + delimiter
	return upperBound + header + strings.Join(usersStr, delimiter) + lowerBound
}

func bookFormat(books []string) {
	for i, book := range books {
		books[i] = fmt.Sprintf("%*s", maxBookLength, book)
	}
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
			averageBookAge[book] += int(user.Age)
		}
	}
	for book := range averageBookAge {
		averageBookAge[book] /= bookCounter[book]
	}
	return averageBookAge

}

func strAvetageAge(averageBookAge map[string]int) string {
	averageTitle := "Avegare age"
	stringAverageAge := fmt.Sprintf("%*s ┃ %s\n", maxBookLength, "Book", averageTitle)
	delimiter := strings.Repeat("━", maxBookLength) + "━╋━" + strings.Repeat("━", len(averageTitle)) + "\n"
	stringAverageAge += delimiter
	for book := range averageBookAge {
		stringAverageAge += fmt.Sprintf("%*s ┃ %*d\n", maxBookLength, book, len(averageTitle), averageBookAge[book])
		stringAverageAge += delimiter
	}
	return stringAverageAge
}

func main() {
	reader := Reader()
	users, _ := Users(reader)
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
