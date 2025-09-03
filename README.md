# mdai

Markdownファイルの内容をAIに質問し、回答を自動で追記するCLIツールです。

## 🚀 機能

- **AI質問**: Markdownファイルの引用部分を抽出してAIに質問
- **自動追記**: AIの回答を元のファイルに自動で追記
- **コスト計算**: OpenAI APIの使用コストを自動計算
- **AI質問**: OpenAIのGPTモデルを使用して質問に回答
- **クロスプラットフォーム**: Windows、macOS、Linuxで動作

## 📋 前提条件

- Go 1.22.0以上
- OpenAI APIキー
    - https://platform.openai.com/api-keys

## 🛠️ インストール

### 方法1: go installを使用（推奨）

```bash
go install github.com/koooyooo/mdai@latest
```

### 方法2: ソースからビルド

#### 1. リポジトリのクローン

```bash
$ git clone https://github.com/koooyooo/mdai.git
$ cd mdai
```

#### 2. 依存関係のインストール

```bash
$ go mod download
```

#### 3. ビルド

```bash
$ go build -o mdai
```

#### 4. 実行可能ファイルをPATHに追加（オプション）

```bash
# macOS/Linux
$ sudo cp mdai /usr/local/bin/

# Windows
# mdai.exeを適切なディレクトリにコピー
```

## 🔑 セットアップ

### OpenAI APIキーの設定

環境変数にOpenAI APIキーを設定してください：

```bash
# macOS/Linux
export OPENAI_API_KEY="your-api-key-here"

# Windows
set OPENAI_API_KEY=your-api-key-here
```

または、`.bashrc`や`.zshrc`に追加して永続化：

```bash
echo 'export OPENAI_API_KEY="your-api-key-here"' >> ~/.bashrc
source ~/.bashrc
```

## 📖 使用方法

### 基本的な使用方法

```bash
# Markdownファイルの引用部分をAIに質問
mdai ask path/to/your/file.md
```

### 使用例

1. **Markdownファイルの準備**

```markdown
# AI学習メモ

> AIを学ぶにあたってコツはありますか？

ここに既存の内容があれば、AIの回答が追記されます。
```

2. **AIに質問**

```bash
mdai ask ai_learning.md
```

3. **結果**

```markdown
# AI学習メモ

> AIを学ぶにあたってコツはありますか？

ここに既存の内容があれば、AIの回答が追記されます。

AIを学ぶにあたってのコツはいくつかあります。まず、基礎知識をしっかりと固めることが重要です...
```

## 💰 コスト計算

mdaiは自動的にAPI使用コストを計算し、ログに表示します。現在使用されているモデルの価格：

- **GPT-4o-mini**: $0.15/1M input, $0.60/1M output（デフォルト）
- **GPT-4o**: $2.50/1M input, $10.00/1M output
- **GPT-4 Turbo**: $10.00/1M input, $30.00/1M output
- **GPT-3.5-turbo**: $0.50/1M input, $1.50/1M output

**注意**: 現在の実装では、GPT-4o-miniがデフォルトモデルとして使用されています。


## 🏗️ プロジェクト構造

```
mdai/
├── cmd/           # CLIコマンド
│   ├── ask.go     # askコマンドの実装
│   └── root.go    # ルートコマンド
├── models/        # AIモデル関連
│   ├── ai_model.go    # AIモデルの定義
│   ├── constants.go    # モデル定数
│   └── helpers.go      # ヘルパー関数
├── util/          # ユーティリティ
│   └── file/      # ファイル操作
├── main.go        # エントリーポイント
└── go.mod         # Goモジュール定義
```

## 🔧 開発

### 依存関係の追加

```bash
go get github.com/package-name
```

### テストの実行

```bash
go test ./...
```

### リントの実行

```bash
# golangci-lintがインストールされている場合
golangci-lint run
```

## 📝 ライセンス

このプロジェクトはMITライセンスの下で公開されています。詳細は[LICENSE](LICENSE)ファイルを参照してください。

## 🤝 コントリビューション

プルリクエストやイシューの報告を歓迎します！

1. このリポジトリをフォーク
2. 機能ブランチを作成 (`git checkout -b feature/amazing-feature`)
3. 変更をコミット (`git commit -m 'Add some amazing feature'`)
4. ブランチにプッシュ (`git push origin feature/amazing-feature`)
5. プルリクエストを作成

### 🚧 開発状況

現在、以下の機能が実装されています：
- OpenAI GPTモデルを使用した質問回答
- Markdownファイルの引用抽出と回答追記
- コスト計算機能

今後の開発予定：
- モデル選択機能の追加
- 他のAIプロバイダー（Claude等）への対応
- 設定ファイルによるカスタマイズ

**注意**: このツールを使用する際は、OpenAI APIの利用規約と料金体系を確認してください。

