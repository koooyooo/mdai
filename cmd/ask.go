/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// askCmd represents the ask command
var askCmd = &cobra.Command{
	Use:   "ask",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ask(args)
	},
}

func init() {
	rootCmd.AddCommand(askCmd)
}

func ask(args []string) error {
	path := args[0]
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	fmt.Println(string(b))

	// INSERT_YOUR_CODE

	// 1. マークダウンから最後の引用部分（> で始まる行）を抽出
	var lastQuote string
	lines := []byte(string(b))
	start := 0
	for i, c := range lines {
		if c == '\n' {
			line := string(lines[start:i])
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
	if start < len(lines) {
		line := string(lines[start:])
		if len(line) > 0 && line[0] == '>' {
			lastQuote = line[1:]
			if len(lastQuote) > 0 && lastQuote[0] == ' ' {
				lastQuote = lastQuote[1:]
			}
		}
	}

	if lastQuote == "" {
		fmt.Println("引用部分（> で始まる行）が見つかりませんでした。")
		return nil
	}

	// 2. OpenAI API キー取得
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("環境変数 OPENAI_API_KEY が設定されていません")
	}

	// 3. OpenAI API へ問い合わせ
	reqBody := fmt.Sprintf(`{
		"model": "gpt-3.5-turbo",
		"messages": [
			{"role": "user", "content": "%s"}
		]
	}`, lastQuote)

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions",
		strings.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("OpenAI API error: %s", string(body))
	}

	type Choice struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	}
	type OpenAIResponse struct {
		Choices []Choice `json:"choices"`
	}
	var openaiResp OpenAIResponse
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&openaiResp); err != nil {
		return err
	}
	if len(openaiResp.Choices) == 0 {
		return fmt.Errorf("OpenAI API から応答がありません")
	}
	answer := openaiResp.Choices[0].Message.Content

	// 4. マークダウンファイルの末尾に結果を追記
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	appendText := fmt.Sprintf("\n\n%s\n", answer)
	if _, err := f.WriteString(appendText); err != nil {
		return err
	}

	fmt.Println("OpenAI の応答をファイル末尾に追記しました。")

	return nil
}
