package utils

import (
	"fmt"
	"time"
)

func generateRefPrefix(timestamp time.Time) string {

	year := timestamp.Year()
	month := timestamp.Month()
	day := timestamp.Day()

	// Determine the year character ('A' for the current year)
	yearCharacter := 'A' + rune((year-2023)%26)

	// Determine the month character ('A' for January, 'B' for February, ..., 'L' for December)
	monthCharacter := 'A' + rune(month-1)

	// Determine the day character ('1' to '9' for the first 9 days, 'A' to 'Z' for the rest)
	var dayCharacter rune
	if day >= 1 && day <= 9 {
		dayCharacter = '0' + rune(day)
	} else {
		dayCharacter = 'A' + rune(day-10)
	}

	// Combine the characters into a string
	result := string([]rune{dayCharacter, monthCharacter, yearCharacter})

	return result
}

func formatTimeComponent(t int) string {

	formatted := fmt.Sprintf("%02d", t)

	return formatted
}

func generateTimeString() string {

	now := time.Now()
	hour := formatTimeComponent(now.Hour())
	minute := formatTimeComponent(now.Minute())
	second := formatTimeComponent(now.Second())
	millisecond := formatTimeComponent(now.Nanosecond() / 1000000)

	return hour + minute + second + millisecond
}

func GenerateTransactionRef() string {
	prefix := generateRefPrefix(time.Now())
	rest := generateTimeString()

	ref := fmt.Sprintf("%s%s", prefix, rest)

	return ref
}
