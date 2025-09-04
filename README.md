# mdai

This is a CLI tool that allows you to ask questions about the contents of a Markdown file to an AI and automatically append the answers.

## 🚀 Features

- **AI Questions**: Extracts the last quoted part from the Markdown file to ask the AI and automatically appends the answer to the original file
- **AI Summary**: Summarizes Markdown files (saves to a separate file)
- **AI Translation**: Translates Markdown files to a specified language (saves to a separate file)
- **Cost Calculation**: Automatically calculates the usage cost of the OpenAI API
- **Cross-Platform**: Works on Windows, macOS, and Linux

## 📋 Prerequisites

- Go 1.22.0 or higher
- OpenAI API key
    - **Required**: Set `OPENAI_API_KEY` environment variable
    - see: https://platform.openai.com/api-keys

For detailed installation and setup instructions, please refer to [INSTALL.md](INSTALL.md).

## 📖 Usage

### Basic Usage

```bash
# Initialize the configuration file (only on first use)
mdai init

# Ask the AI about quoted parts of the Markdown file
mdai answer path/to/your/file.md

# Summarize the contents of the Markdown file
mdai summarize path/to/your/file.md

# Translate the Markdown file to a specified language
mdai translate path/to/your/file.md ja
```

### Customizing the Configuration File

mdai can use a configuration file to customize its operation. The configuration file is located at `~/.mdai/config.yml`.

#### Initializing the Configuration File

```bash
# Initialize the configuration file (first setup)
mdai init
```

This command performs the following actions:
1. Creates the `~/.mdai` directory
2. Copies `config.sample.yml` to `~/.mdai/config.yml`
3. Displays the path to the configuration file

#### Configuration Items

- **Default Settings**: AI model, quality settings, log level
- **answer Command**: System message, target character count
- **summarize Command**: System message, target character count
- **translate Command**: System message

For detailed configuration examples, refer to `config/config.sample.yml`.

### Usage Example

1. **Prepare the Markdown File**

```markdown
# AI Learning Notes

> Are there any tips for learning AI?

If there is existing content here, the AI's answer will be appended.
```

2. **Ask the AI**

```bash
mdai answer ai_learning.md
```

3. **Result**

```markdown
# AI Learning Notes

> Are there any tips for learning AI?

If there is existing content here, the AI's answer will be appended.

There are several tips for learning AI. First, it is important to solidify your foundational knowledge...
```

### Translation Example

```bash
# Translate to English
mdai translate ai_learning.md en

# Translate to Japanese
mdai translate ai_learning.md ja
```

The translation results will be saved as `ai_learning_en.md` and `ai_learning_ja.md`.

## 💰 Cost Calculation

mdai automatically calculates API usage costs and displays them in the logs.

**Note**: Currently, the default model being used is GPT-4o-mini. Please check the [OpenAI pricing page](https://openai.com/pricing) for current model prices.

## 🏗️ Project Structure

```
mdai/
├── cmd/           # CLI commands
│   ├── answer.go     # Implementation of the answer command
│   ├── summarize.go  # Implementation of the summarize command
│   ├── translate.go  # Implementation of the translate command
│   ├── init.go       # Implementation of the init command
│   └── root.go       # Root command
├── config/        # Configuration files
│   └── config.go     # Configuration struct and loading process
├── config.sample.yml # Sample configuration file
├── controller/    # AI control
│   └── controller.go # OpenAI API control
├── models/        # AI model related
│   ├── ai_model.go    # Definition of AI models
│   ├── constants.go    # Model constants
│   └── helpers.go      # Helper functions
├── util/          # Utilities
│   └── file/      # File operations
├── mdai.go        # Entry point
└── go.mod         # Go module definition
```

## 🔧 Development

### Adding Dependencies

```bash
go get github.com/package-name
```

### Running Tests

```bash
go test ./...
```

### Running Lint

```bash
# If golangci-lint is installed
golangci-lint run
```

## 📝 License

This project is licensed under the MIT License. Please refer to the [LICENSE](LICENSE) file for details.

## 🤝 Contribution

Pull requests and issue reports are welcome!

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Create a pull request

### 🚧 Development Status

Currently, the following features have been implemented:
- Question answering using OpenAI GPT models
- Extraction of quoted parts from Markdown files and appending answers
- Cost calculation feature

Planned developments include:
- Adding a model selection feature
- Support for other AI providers (e.g., Claude)
- Customization through configuration files

**Note**: Please check the OpenAI API terms of service and pricing structure when using this tool.

## 🔗 Related Links

- [INSTALL.md](INSTALL.md) - Installation and setup instructions
- [LICENSE](LICENSE) - License information