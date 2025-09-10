/*
Copyright © 2025 koooyooo
*/
package file

import (
	"runtime"
	"strings"
)

func getLineSeparator() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func LoadLastQuote(content string) (bool, string, string, error) {
	lineSep := getLineSeparator()
	lines := strings.Split(content, lineSep)

	var otherContents string
	var lastQuoteLines []string
	var currentQuoteLines []string
	var inQuote = false

	// 行を順番に処理して、最後の有効な引用を特定
	for _, line := range lines {
		if strings.HasPrefix(line, ">") {
			// 引用行の場合
			if !inQuote {
				// 新しい引用ブロックの開始
				// 前の有効な引用がある場合は、それをotherContentsに移動
				if len(lastQuoteLines) > 0 {
					for _, quoteLine := range lastQuoteLines {
						otherContents += quoteLine + lineSep
					}
					lastQuoteLines = nil
				}
				currentQuoteLines = []string{line}
				inQuote = true
			} else {
				// 既存の引用ブロックの継続
				currentQuoteLines = append(currentQuoteLines, line)
			}
		} else {
			// 非引用行の場合
			if inQuote {
				// 引用ブロックが終了
				// この行が空行またはスペースのみかチェック
				trimmed := strings.TrimSpace(line)
				if trimmed == "" {
					// 空行またはスペースのみの場合、一時的に引用を保存
					// ただし、その後に通常のテキストが続く場合は無効化される
					lastQuoteLines = make([]string, len(currentQuoteLines))
					copy(lastQuoteLines, currentQuoteLines)
					currentQuoteLines = nil
					inQuote = false
					otherContents += line + lineSep
				} else {
					// 非空白文字がある場合、現在の引用をクリア（無効とする）
					// この引用ブロックは最後の引用ではない
					// 無効化された引用部分もotherContentsに含める
					for _, quoteLine := range currentQuoteLines {
						otherContents += quoteLine + lineSep
					}
					currentQuoteLines = nil
					inQuote = false
					otherContents += line + lineSep
					// 前の引用も無効化
					lastQuoteLines = nil
				}
			} else {
				// 引用ブロック外の通常行
				otherContents += line + lineSep
				// 通常のテキストが続いている場合、前の引用を無効化
				trimmed := strings.TrimSpace(line)
				if trimmed != "" {
					lastQuoteLines = nil
				}
			}
		}
	}

	// ファイル末尾で引用ブロックが終了している場合
	if inQuote {
		// ファイル末尾の引用は有効な最後の引用として扱う
		lastQuoteLines = make([]string, len(currentQuoteLines))
		copy(lastQuoteLines, currentQuoteLines)
	}

	// 有効な引用が見つからない場合
	if len(lastQuoteLines) == 0 {
		return false, "", otherContents, nil
	}

	// 引用内容を構築
	var lastQuote string
	for i, line := range lastQuoteLines {
		if i > 0 {
			lastQuote += lineSep
		}
		// '>' の後の部分を取得
		quoteContent := line[1:]
		// 先頭のスペースを削除
		if len(quoteContent) > 0 && quoteContent[0] == ' ' {
			quoteContent = quoteContent[1:]
		}
		lastQuote += quoteContent
	}

	return true, lastQuote, otherContents, nil
}
