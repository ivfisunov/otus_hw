package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	runeString := []rune(str)
	var runeStringLengtht = len(runeString)
	var result = ""
	for i := 0; i < runeStringLengtht; {
		letter := runeString[i]
		if unicode.IsDigit(letter) {
			return "", ErrInvalidString
		}

		// escape symbol implementation
		if string(letter) == "\\" {
			escapedSymbol := runeString[i+1]
			var countOfEscapedSymbol rune
			if i < (runeStringLengtht - 2) { // prevent out of range
				countOfEscapedSymbol = runeString[i+2]
			}
			if unicode.IsDigit(countOfEscapedSymbol) {
				repeatCount, _ := strconv.Atoi(string(countOfEscapedSymbol))
				result += string(strings.Repeat(string(escapedSymbol), repeatCount))
				i += 3
				continue
			}
			result += string(escapedSymbol)
			i += 2
			continue
		}

		var number rune
		if i < (runeStringLengtht - 1) { // prevent out of range
			number = runeString[i+1]
		}
		if unicode.IsDigit(number) {
			repeatCount, _ := strconv.Atoi(string(number))
			result += string(strings.Repeat(string(letter), repeatCount))
			i += 2
			continue
		}
		result += string(letter)
		i++
	}

	return result, nil
}
