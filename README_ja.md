# mdai

Markdownファイルの内容をAIに質問し、回答を自動で追記するCLIツールです。

## 🚀 機能

- **AI質問**: Markdownファイルの引用部分を抽出してAIに質問
- **自動追記**: AIの回答を元のファイルに自動で追記
- **コスト計算**: OpenAI APIの使用コストを自動計算
- **AI質問**: OpenAIのGPTモデルを使用して質問に回答
- **AI翻訳**: Markdownファイルを指定言語に翻訳
- **クロスプラットフォーム**: Windows、macOS、Linuxで動作

## 📋 前提条件

- Go 1.22.0以上
- OpenAI APIキー
    - **必須**: `OPENAI_API_KEY`環境変数を設定してください
    - see: https://platform.openai.com/api-keys

詳細なインストールとセットアップ手順は [INSTALL.md](INSTALL.md) を参照してください。

## 📖 使用方法

### 基本的な使用方法

```bash
# 設定ファイルの初期化（初回のみ）
mdai init

# Markdownファイルの引用部分をAIに質問
mdai answer path/to/your/file.md

# Markdownファイルの内容を要約
mdai summarize path/to/your/file.md

# Markdownファイルを指定言語に翻訳
mdai translate path/to/your/file.md ja
```

### 設定ファイルのカスタマイズ

mdaiは設定ファイルを使用して動作をカスタマイズできます。設定ファイルは `~/.mdai/config.yml` に配置されます。

#### 設定ファイルの初期化

```bash
# 設定ファイルを初期化（初回セットアップ）
mdai init
```

このコマンドは以下を実行します：
1. `~/.mdai` ディレクトリを作成
2. `config.sample.yml` を `~/.mdai/config.yml` にコピー
3. 設定ファイルのパスを表示

#### 設定項目

- **デフォルト設定**: AIモデル、品質設定、ログレベル
- **answerコマンド**: システムメッセージ、目標文字数
- **summarizeコマンド**: システムメッセージ、目標文字数
- **translateコマンド**: システムメッセージ

詳細な設定例は `config/config.sample.yml` を参照してください。

### 使用例

1. **Markdownファイルの準備**

```markdown
# AI学習メモ

> AIを学ぶにあたってコツはありますか？

ここに既存の内容があれば、AIの回答が追記されます。
```

2. **AIに質問**

```bash
mdai answer ai_learning.md
```

3. **結果**

```markdown
# AI学習メモ

> AIを学ぶにあたってコツはありますか？

ここに既存の内容があれば、AIの回答が追記されます。

AIを学ぶにあたってのコツはいくつかあります。まず、基礎知識をしっかりと固めることが重要です...

### 翻訳の例

```bash
# 英語に翻訳
mdai translate ai_learning.md en

# 日本語に翻訳
mdai translate ai_learning.md ja
```

翻訳結果は `ai_learning_en.md`、`ai_learning_ja.md` として保存されます。
```

## 💰 コスト計算

mdaiは自動的にAPI使用コストを計算し、ログに表示します。

**注意**: 現在の実装では、GPT-4o-miniがデフォルトモデルとして使用されています。現在のモデル価格については[OpenAI料金ページ](https://openai.com/pricing)をご確認ください。


## 🏗️ プロジェクト構造

```
mdai/
├── cmd/           # CLIコマンド
│   ├── answer.go     # answerコマンドの実装
│   ├── summarize.go  # summarizeコマンドの実装
│   ├── translate.go  # translateコマンドの実装
│   ├── init.go       # initコマンドの実装
│   └── root.go       # ルートコマンド
├── config/        # 設定ファイル
│   └── config.go     # 設定構造体と読み込み処理
├── config.sample.yml # サンプル設定ファイル
├── controller/    # AI制御
│   └── controller.go # OpenAI API制御
├── models/        # AIモデル関連
│   ├── ai_model.go    # AIモデルの定義
│   ├── constants.go    # モデル定数
│   └── helpers.go      # ヘルパー関数
├── util/          # ユーティリティ
│   └── file/      # ファイル操作
├── mdai.go        # エントリーポイント
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

## 🔗 関連リンク

- [INSTALL.md](INSTALL.md) - インストールとセットアップ手順
- [LICENSE](LICENSE) - ライセンス情報

