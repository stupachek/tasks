package format

import (
	"encoding/base64"
	"epam/model"
	"fmt"
	"strings"
)

const (
	nameLength = 8
	Name       = "Name"
	Books      = "Books"
	Age        = "Age"
	Active     = "Active"
	Mass       = "Mass"
)

var (
	maxNameLength   = len(Name)
	maxBookLength   = len(Books)
	maxAgeLength    = len(Age)
	maxActiveLenght = len(Active)
	maxMassLength   = len(Mass)
)

func DelimiterPattern(vertical, horizontal string) string {
	return strings.Join([]string{strings.Repeat(vertical, maxNameLength), strings.Repeat(vertical, maxAgeLength), strings.Repeat(vertical, maxActiveLenght), strings.Repeat(vertical, maxMassLength-2), strings.Repeat(vertical, maxBookLength)}, horizontal)
}

func FormatUsers(users []model.User) string {
	usersStr := make([]string, len(users))
	upperBound := "┏" + DelimiterPattern("━", "━┳━") + "┓\n"
	spacesBeforeBook := DelimiterPattern(" ", " ┃ ")
	spacesBeforeBook = spacesBeforeBook[:len(spacesBeforeBook)-maxBookLength]
	lowerBound := "\n┗" + DelimiterPattern("━", "━┻━") + "┛"
	for i, user := range users {
		name := fmt.Sprintf("%*s", maxNameLength, users[i].Name)
		age := fmt.Sprintf("%*d", maxAgeLength, user.Age)
		activestr := fmt.Sprintf("%*s", maxActiveLenght, "-")
		if user.Active {
			activestr = fmt.Sprintf("%*s", maxActiveLenght, "+")
		}
		mass := fmt.Sprintf("% *.2f", maxMassLength-2, user.Mass)
		BookFormat(user.Books)
		books := strings.Join(users[i].Books, "┃\n┃"+spacesBeforeBook)
		usersStr[i] = "┃" + strings.Join([]string{name, age, activestr, mass, books}, " ┃ ") + "┃"
	}
	delimiter := "\n┣" + DelimiterPattern("━", "━╋━") + "┫\n"
	header := fmt.Sprintf("┃%*s ┃ %*s ┃ %*s ┃ %*s ┃ %*s┃", maxNameLength, Name, maxAgeLength, Age, maxActiveLenght, Active, maxMassLength-2, Mass, maxBookLength, Books) + delimiter
	return upperBound + header + strings.Join(usersStr, delimiter) + lowerBound
}

func BookFormat(books []string) {
	for i, book := range books {
		books[i] = fmt.Sprintf("%*s", maxBookLength, book)
	}
}

func PrepareUsers(users []model.User) {
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
func count(n uint8) int {
	counter := 0
	for i := 0; i < 8; i++ {
		counter += int(n & 1)
		n >>= 1
	}
	return counter
}

func averageAge(users []model.User) map[string]int {
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
