# Installation and Setup

This document provides detailed instructions on installing and setting up mdai.

## 📋 Prerequisites

- Go version 1.22.0 or higher
- OpenAI API key
    - see: https://platform.openai.com/api-keys

## 🛠️ Installation

### Method 1: Using `go install` (Recommended)

```bash
go install github.com/koooyooo/mdai@latest
```

### Method 2: Building from Source

#### 1. Clone the Repository

```bash
$ git clone https://github.com/koooyooo/mdai.git
$ cd mdai
```

#### 2. Install Dependencies

```bash
$ go mod download
```

#### 3. Build

```bash
$ go build -o mdai
```

#### 4. (Optional) Add the Executable to PATH

```bash
# macOS/Linux
$ sudo cp mdai /usr/local/bin/

# Windows
# Copy mdai.exe to an appropriate directory
```

## 🔑 Setup

### Setting the OpenAI API Key

Please set the OpenAI API key in your environment variables:

```bash
# macOS/Linux
export OPENAI_API_KEY="your-api-key-here"

# Windows
set OPENAI_API_KEY=your-api-key-here
```

Alternatively, you can add it to `.bashrc` or `.zshrc` for persistence:

```bash
echo 'export OPENAI_API_KEY="your-api-key-here"' >> ~/.bashrc
source ~/.bashrc
```

## ✅ Verifying the Installation

To check if the installation was successful:

```bash
mdai --version
```

Or

```bash
mdai --help
```

## 🚨 Troubleshooting

### Common Issues

1. **Command not found**
   - Ensure the Go PATH is set correctly
   - Check the Go path using `go env GOPATH`

2. **API key not recognized**
   - Ensure the environment variable is set correctly
   - Restart the terminal to reload the environment variable

3. **Permission errors**
   - Check if the executable has execution permissions
   - Grant permissions using `chmod +x mdai`

## 🔗 Related Links

- [README.md](README.md) - Project overview and usage instructions
- [LICENSE](LICENSE) - License information