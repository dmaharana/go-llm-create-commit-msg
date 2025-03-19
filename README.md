# createCommitMsg

This project provides a command-line tool to generate commit messages.

## Prerequisites

- Go (version 1.18 or higher)
- Git

## Installation

1.  Clone the repository:

    ```bash
    git clone https://github.com/titu/createCommitMsg
    cd createCommitMsg
    ```

    (Replace `<repository_url>` with the actual URL of the repository.)

2.  Build the executable:

    ```bash
    make build
    ```

## Usage

1.  Set your OpenRouter API key as an environment variable:

    ```bash
    export OPEN_ROUTER_KEY="YOUR_API_KEY"
    ```

    (Replace `YOUR_API_KEY` with your actual OpenAI API key.)

2.  Run the executable:

    ```bash
    ./bin/createCommitMsg
    ```

    or

    ```bash
    make
    ```

    And then run the script after setting the OpenRouter API key.

    This will generate a commit message based on the staged changes in your Git repository.

## Contributing

We welcome contributions! Please follow these guidelines:

1.  Fork the repository.
2.  Create a new branch for your feature or bug fix.
3.  Make your changes and commit them with clear, concise commit messages.
4.  Ensure your code passes all tests.
5.  Submit a pull request.

## Project Structure

- `cmd/`: Contains the main application logic.
- `internal/`: Contains internal packages.
  - `internal/git/`: Git related functionalities.
  - `internal/llm/`: LLM (Language Model) related functionalities.
- `scripts/`: Contains shell scripts.
  - `scripts/createMsg.sh`: Shell script to create commit message.
- `Makefile`: Contains build instructions.
- `go.mod`: Go module file.
- `go.sum`: Go module checksum file.
- `.gitignore`: Specifies intentionally untracked files that Git should ignore.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
