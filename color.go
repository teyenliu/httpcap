package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

const (
	Gray = uint8(iota + 90)
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	EndColor         = "\033[0m"
	printV           = "A"
	contentJsonRegex = `application/json`
)

func Color(str string, color uint8) string {
	return fmt.Sprintf("%s%s%s", ColorStart(color), str, EndColor)
}

func ColorStart(color uint8) string {
	return fmt.Sprintf("\033[%dm", color)
}

func ColorfulRequest(str string) string {
	idx := 0
	lines := strings.Split(str, "\n")
	if printV == "A" || printV == "H" {
		strs := strings.Split(lines[0], " ")
		if len(strs) >= 3 {
			strs[0] = Color(strs[0], Magenta)
			strs[1] = Color(strs[1], Cyan)
			strs[2] = Color(strs[2], Magenta)
			lines[0] = strings.Join(strs, " ")
			idx = 1
		}
	}
	for i, line := range lines[idx:] {
		substr := strings.Split(line, ":")
		if len(substr) < 2 {
			continue
		}
		substr[0] = Color(substr[0], Gray)
		substr[1] = Color(strings.Join(substr[1:], ":"), Cyan)
		lines[i+1] = strings.Join(substr[:2], ":")
	}
	return strings.Join(lines, "\n")
}

func ColorfulResponse(str string) string {
	if isJSON(str) {
		str = ColorfulJson(str)
	} else {
		str = ColorfulHTML(str)
	}
	return str
}

func ColorfulJson(str string) string {
	var rsli []rune
	var key, val, startcolor, endcolor, startsemicolon bool
	var prev rune
	for _, char := range []rune(str) {
		switch char {
		case ' ':
			rsli = append(rsli, char)
		case '{':
			startcolor = true
			key = true
			val = false
			rsli = append(rsli, char)
		case '}':
			startcolor = false
			endcolor = false
			key = false
			val = false
			rsli = append(rsli, char)
		case '"':
			if startsemicolon && prev == '\\' {
				rsli = append(rsli, char)
			} else {
				if startcolor {
					rsli = append(rsli, char)
					if key {
						rsli = append(rsli, []rune(ColorStart(Magenta))...)
					} else if val {
						rsli = append(rsli, []rune(ColorStart(Cyan))...)
					}
					startsemicolon = true
					key = false
					val = false
					startcolor = false
				} else {
					rsli = append(rsli, []rune(EndColor)...)
					rsli = append(rsli, char)
					endcolor = true
					startsemicolon = false
				}
			}
		case ',':
			if !startsemicolon {
				startcolor = true
				key = true
				val = false
				if !endcolor {
					rsli = append(rsli, []rune(EndColor)...)
					endcolor = true
				}
			}
			rsli = append(rsli, char)
		case ':':
			if !startsemicolon {
				key = false
				val = true
				startcolor = true
				if !endcolor {
					rsli = append(rsli, []rune(EndColor)...)
					endcolor = true
				}
			}
			rsli = append(rsli, char)
		case '\n', '\r', '[', ']':
			rsli = append(rsli, char)
		default:
			if !startsemicolon {
				if key && startcolor {
					rsli = append(rsli, []rune(ColorStart(Magenta))...)
					key = false
					startcolor = false
					endcolor = false
				}
				if val && startcolor {
					rsli = append(rsli, []rune(ColorStart(Cyan))...)
					val = false
					startcolor = false
					endcolor = false
				}
			}
			rsli = append(rsli, char)
		}
		prev = char
	}
	return string(rsli)
}

func ColorfulHTML(str string) string {
	return Color(str, Green)
}

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil

}
