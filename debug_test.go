package main

import (
	"fmt"
	"os"

	"github.com/koooyooo/mdai/util/file"
)

func main() {
	// Sample.mdの内容を読み込む
	content, err := os.ReadFile("Sample.md")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	success, quote, other, err := file.LoadLastQuote(string(content))
	fmt.Printf("Success: %v\n", success)
	fmt.Printf("Quote: %q\n", quote)
	fmt.Printf("Other length: %d\n", len(other))
	fmt.Printf("Error: %v\n", err)
}
