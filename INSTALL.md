# インストールとセットアップ

このドキュメントでは、mdaiのインストールとセットアップについて詳しく説明します。

## 📋 前提条件

- Go 1.22.0以上
- OpenAI APIキー
    - see: https://platform.openai.com/api-keys

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

## ✅ インストール確認

インストールが完了したかどうかを確認するには：

```bash
mdai --version
```

または

```bash
mdai --help
```

## 🚨 トラブルシューティング

### よくある問題

1. **コマンドが見つからない**
   - GoのPATHが正しく設定されているか確認
   - `go env GOPATH`でGoのパスを確認

2. **APIキーが認識されない**
   - 環境変数が正しく設定されているか確認
   - ターミナルを再起動して環境変数を再読み込み

3. **権限エラー**
   - 実行可能ファイルに実行権限があるか確認
   - `chmod +x mdai`で権限を付与

## 🔗 関連リンク

- [README.md](README.md) - プロジェクト概要と使用方法
- [LICENSE](LICENSE) - ライセンス情報
