# createCommitMsg

A command-line tool that uses Large Language Models (LLMs) via OpenRouter to generate **commit messages** and **code review comments** based on your staged Git changes.

---

## Prerequisites

- Go (version 1.18 or higher)
- Git
- An OpenRouter API key (https://openrouter.ai/)

---

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/titu/createCommitMsg
   cd createCommitMsg
   ```

2. **Build the executable:**

   ```bash
   make build
   ```

---

## Configuration

Set your OpenRouter API key as an environment variable:

```bash
export OPEN_ROUTER_KEY="YOUR_API_KEY"
```

Optionally, specify a custom LLM model (default is `google/gemini-2.0-flash-lite-preview-02-05:free`):

```bash
export MODEL_NAME="your-preferred-model-name"
```

---

## Usage

### Using the provided script

Run the helper script with optional flags:

```bash
scripts/createMsg.sh [-r | -c | -b]
```

- `-r` : Generate **code review comments** only
- `-c` : Generate **commit message** only
- `-b` : Generate **both** (default if no flag is provided)

### Running directly with Go

Alternatively, run the tool directly with customizable flags:

```bash
go run cmd/main.go --mode [r|c|b] --output [r|c|b] --format [r|j]
```

- `--mode`:

  - `r` or `review` : generate **code review comments** only
  - `c` or `comment` : generate **commit message** only
  - `b` or `both` (default) : generate **both**

- `--output`:

  - `r` or `review` : display **code review comments** only
  - `c` or `comment` : display **commit message** only
  - `b` or `both` (default) : display **both**

- `--format`:
  - `r` or `raw` (default) : plain text output
  - `j` or `json` : JSON formatted output

---

## How it works

- The tool analyzes your **staged Git changes**.
- It sends the changes to an LLM via OpenRouter.
- Depending on the selected mode, it generates:
  - A **commit message**
  - **Code review comments**
  - Or both
- The output is displayed in your terminal, either as plain text or JSON.

---

## Project Structure

- `cmd/` — Main CLI application
- `internal/`
  - `action/` — Action constants
  - `constant/` — Constant values
  - `display/` — Display helpers
  - `git/` — Git integration
  - `llm/` — LLM API integration
  - `prompt/` — Prompt templates
- `scripts/`
  - `createMsg.sh` — Helper script to run the tool
- `Makefile` — Build commands
- `go.mod`, `go.sum` — Go module files
- `.gitignore` — Git ignore rules

---

## Contributing

We welcome contributions! Please:

1. Fork the repository
2. Create a feature branch
3. Make your changes with clear commit messages
4. Ensure all tests pass
5. Submit a pull request

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
