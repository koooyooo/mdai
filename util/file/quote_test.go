/*
Copyright © 2025 koooyooo
*/
package file

import (
	"testing"
)

func TestLoadLastQuote(t *testing.T) {
	tests := []struct {
		name            string
		content         string
		expectedSuccess bool
		expectedQuote   string
		expectedOther   string
		expectedError   string
	}{
		{
			name: "単一の引用ブロック（ファイル末尾）",
			content: `通常のテキスト
> これは引用です
> 複数行の引用`,
			expectedSuccess: true,
			expectedQuote:   "これは引用です\n複数行の引用",
			expectedOther:   "通常のテキスト\n",
			expectedError:   "",
		},
		{
			name: "引用ブロック後に非空白文字がある場合（無効）",
			content: `> 最初の引用
> 続き
通常のテキスト

> 2番目の引用
> 続き

> 3番目の引用`,
			expectedSuccess: true,
			expectedQuote:   "3番目の引用",
			expectedOther:   "> 最初の引用\n> 続き\n通常のテキスト\n\n\n> 2番目の引用\n> 続き\n",
			expectedError:   "",
		},
		{
			name: "引用ブロック後に空行のみ（有効）",
			content: `通常のテキスト

> 引用ブロック
> 複数行

`,
			expectedSuccess: true,
			expectedQuote:   "引用ブロック\n複数行",
			expectedOther:   "通常のテキスト\n\n\n\n",
			expectedError:   "",
		},
		{
			name: "引用ブロック後にスペースのみ（有効）",
			content: `通常のテキスト

> 引用ブロック
> 複数行
    `,
			expectedSuccess: true,
			expectedQuote:   "引用ブロック\n複数行",
			expectedOther:   "通常のテキスト\n\n    \n",
			expectedError:   "",
		},
		{
			name: "引用が見つからない場合",
			content: `通常のテキスト
別のテキスト`,
			expectedSuccess: false,
			expectedQuote:   "",
			expectedOther:   "通常のテキスト\n別のテキスト\n",
			expectedError:   "",
		},
		{
			name: "引用ブロック後に非空白文字がある場合（複数パターン）",
			content: `> 1番目の引用
テキスト1

> 2番目の引用
テキスト2

> 3番目の引用`,
			expectedSuccess: true,
			expectedQuote:   "3番目の引用",
			expectedOther:   "> 1番目の引用\nテキスト1\n\n> 2番目の引用\nテキスト2\n\n",
			expectedError:   "",
		},
		{
			name: "引用ブロック内にスペース付きの引用",
			content: `通常のテキスト

> これは引用です
>    スペース付きの引用
> 通常の引用`,
			expectedSuccess: true,
			expectedQuote:   "これは引用です\n   スペース付きの引用\n通常の引用",
			expectedOther:   "通常のテキスト\n\n",
			expectedError:   "",
		},
		{
			name:            "空のファイル",
			content:         "",
			expectedSuccess: false,
			expectedQuote:   "",
			expectedOther:   "\n",
			expectedError:   "",
		},
		{
			name: "引用のみのファイル",
			content: `> 引用のみ
> 複数行`,
			expectedSuccess: true,
			expectedQuote:   "引用のみ\n複数行",
			expectedOther:   "",
			expectedError:   "",
		},
		{
			name: "引用ブロック後に全角スペースのみ（有効）",
			content: `通常のテキスト

> 引用ブロック
　　`,
			expectedSuccess: true,
			expectedQuote:   "引用ブロック",
			expectedOther:   "通常のテキスト\n\n　　\n",
			expectedError:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			success, quote, other, err := LoadLastQuote(tt.content)

			// 成功フラグの確認
			if success != tt.expectedSuccess {
				t.Errorf("LoadLastQuote() success = %v, want %v", success, tt.expectedSuccess)
			}

			// 引用内容の確認
			if quote != tt.expectedQuote {
				t.Errorf("LoadLastQuote() quote = %q, want %q", quote, tt.expectedQuote)
			}

			// その他の内容の確認
			if other != tt.expectedOther {
				t.Errorf("LoadLastQuote() other = %q, want %q", other, tt.expectedOther)
			}

			// エラーの確認
			if tt.expectedError != "" {
				if err == nil {
					t.Errorf("LoadLastQuote() expected error %q, got nil", tt.expectedError)
				} else if err.Error() != tt.expectedError {
					t.Errorf("LoadLastQuote() error = %q, want %q", err.Error(), tt.expectedError)
				}
			} else {
				if err != nil {
					t.Errorf("LoadLastQuote() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestLoadLastQuote_EdgeCases(t *testing.T) {
	t.Run("引用ブロックが複数回無効化される場合", func(t *testing.T) {
		content := `> 1番目
テキスト1

> 2番目
テキスト2

> 3番目
テキスト3

> 最後の引用`

		success, quote, _, err := LoadLastQuote(content)

		if !success {
			t.Errorf("Expected success, got false")
		}
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if quote != "最後の引用" {
			t.Errorf("Expected '最後の引用', got %q", quote)
		}
	})

	t.Run("引用ブロック後に改行とスペースが混在", func(t *testing.T) {
		content := `通常のテキスト

> 引用ブロック
> 複数行

   
   全角スペース
   
`

		success, quote, _, err := LoadLastQuote(content)

		if success {
			t.Errorf("Expected failure (quote followed by actual text), got success")
		}
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		if quote != "" {
			t.Errorf("Expected empty quote, got %q", quote)
		}
	})
}
