/*
Copyright © 2025 koooyooo
*/
package file

import (
	"fmt"
	"runtime"
	"strings"
)

func getLineSeparator() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func LoadLastQuote(content string) (string, string, error) {
	lineSep := getLineSeparator()
	lines := strings.Split(content, lineSep)
	var lastQuote string
	var otherContents string

	for _, line := range lines {
		if !strings.HasPrefix(line, ">") {
			otherContents += line + lineSep
			continue
		}
		if len(line) > 0 && line[0] == '>' {
			lastQuote = line[1:] // '>' の後ろの部分
			if len(lastQuote) > 0 && lastQuote[0] == ' ' {
				lastQuote = lastQuote[1:]
			}
		}
	}

	if lastQuote == "" {
		return "", "", fmt.Errorf("no quote (line starting with >) found")
	}
	return lastQuote, otherContents, nil
}
