package file

import (
	"fmt"
)

func LoadLastQuote(content string) (string, error) {
	// 1. マークダウンから最後の引用部分（> で始まる行）を抽出
	var lastQuote string
	start := 0
	for i, c := range content {
		if c == '\n' {
			line := string(content[start:i])
			if len(line) > 0 && line[0] == '>' {
				lastQuote = line[1:] // '>' の後ろの部分
				if len(lastQuote) > 0 && lastQuote[0] == ' ' {
					lastQuote = lastQuote[1:]
				}
			}
			start = i + 1
		}
	}
	// ファイルの最後の行が改行で終わっていない場合の対応
	if start < len(content) {
		line := string(content[start:])
		if len(line) > 0 && line[0] == '>' {
			lastQuote = line[1:]
			if len(lastQuote) > 0 && lastQuote[0] == ' ' {
				lastQuote = lastQuote[1:]
			}
		}
	}

	if lastQuote == "" {
		fmt.Println("引用部分（> で始まる行）が見つかりませんでした。")
		return "", fmt.Errorf("引用部分（> で始まる行）が見つかりませんでした。")
	}
	return lastQuote, nil
}
