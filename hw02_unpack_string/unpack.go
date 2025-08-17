package hw02unpackstring

import (
	"errors"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	sr := []rune(s)
	var lastSymbol rune
	var s2 string
	var n int
	var backslash bool

	for i, item := range sr {

		// если первый символ в строке цифра
		if unicode.IsDigit(item) && i == 0 {
			return "", ErrInvalidString
		}
		// если символ в строке цифра, перед символом цифра и перед обоими не обратный слеш
		if unicode.IsDigit(item) && unicode.IsDigit(sr[i-1]) && sr[i-2] != '\\' {
			return "", ErrInvalidString
		}
		// если символ обратный слеш и backslash = false (по умолчанию)
		if item == '\\' && !backslash {
			backslash = true
			continue
		}
		// если backslash = true и символ буква
		if backslash && unicode.IsLetter(item) {
			return "", ErrInvalidString
		}
		// если backslash = true
		if backslash {
			s2 += string(item)
			backslash = false
			continue
		}

		// множитель предыдущей буквы
		if unicode.IsDigit(item) {
			n = int(item - '0')
			// если ноль...
			if n == 0 {
				s2 = s2[:len(s2)-len(string(lastSymbol))]
				continue
			}
			// если не ноль...
			for j := 0; j < n-1; j++ {
				s2 += string(sr[i-1])
			}
			lastSymbol = sr[i-1]
			continue
		}
		lastSymbol = item
		s2 += string(item)
	}

	return s2, nil
}
