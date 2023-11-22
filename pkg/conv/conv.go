package conv

import (
	"strconv"
	"strings"
	"unicode"
)

// TODO: check TrimSpace()
func StrtoIntWithoutSpace(text string) (int, error) {
	var number string
	var isLastSymbol bool
	text = strings.ReplaceAll(text, " ", "")
	for _, t := range text {
		if isLastSymbol {
			number = ""
			if unicode.IsDigit(t) {
				number += string(t)
			}
			isLastSymbol = false
		} else if unicode.IsDigit(t) {
			number += string(t)
		} else if string(t) == " " {
			continue
		} else {
			isLastSymbol = true
		}
	}
	return strconv.Atoi(number)
}

func StrToStrWithoutEnter(text string) (string, error) {
	var str string
	for _, t := range text {
		if string(t) != `\n` {
			str += string(t)
		} else {
			return str, nil
		}
	}
	return str, nil
}

func StrToStrLastElement(text string) (string, error) {
	var str string
	var counterYear int
	for _, t := range text {
		if counterYear == 4 {
			str = ""
			counterYear = 0
		}
		if unicode.IsDigit(t) {
			counterYear++
		} else {
			counterYear = 0
		}
		str += string(t)
	}
	return str, nil
}
