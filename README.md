# Fed - AI Agent

> This whole readme is written using this agent

Fed is an AI-driven agent designed to assist with coding, refactoring, and executing bash scripts. It integrates with LLMs (Gemini and Ollama) to understand user prompts, plan tasks, and execute them by interacting with the file system and shell.

## Architecture

The core of the project is the `Agent` struct, which orchestrates tasks using several key interfaces:

- **`AI`**: Interacts with the underlying LLM (Gemini or Ollama) to generate structured JSON responses from prompts.
- **`ToolChain`**: Provides tools for the agent to interact with the environment. It includes functionalities like `RunFileChange` and `ReadFile` to safely read and modify project files.
- **`AgentReporter`**: Handles reporting the agent's status, logging messages, and communicating the results of its operations to the user.

### Integrations

- **Gemini**: Uses the `google.golang.org/genai` client for interacting with Google's Gemini models.
- **Ollama**: Uses a custom HTTP client to interact with local Ollama instances.

## Internal Workflow

1. **Chat Loop**: The agent runs in a continuous loop, waiting for user input via the CLI.
2. **Planning & Thinking**: Upon receiving a prompt, the agent analyzes the context (project files, chat history) and generates a plan or a direct response.
3. **File Editing**: If the task involves code changes, the agent generates a search/replace patch. The `ToolChain` checks file rights and prompts the user for approval before applying any changes.

## Getting Started

### Prerequisites

- An API key for Gemini OR a running local instance of Ollama

### Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
2. Edit `.env` to configure your preferred LLM:
   - For **Ollama**: Set `OLLAMA_ENDPOINT` and `OLLAMA_MODEL`.
   - For **Gemini**: Set `GEMINI_TOKEN` and `GEMINI_MODEL`.

### Installation

There are always pre-built binaries in [releases](https://github.com/hurtki/fed/releases)

#### Oneliners for UNIX-like systems

### Linux (x86_64 / AMD64)

```
curl -L https://github.com/hurtki/fed/releases/latest/download/fed-linux-amd64 -o fed
chmod +x fed
sudo mv fed /usr/local/bin/
```

### macOS Intel

```
curl -L https://github.com/hurtki/fed/releases/latest/download/fed-darwin-amd64 -o fed
chmod +x fed
sudo mv fed /usr/local/bin/
```

### Linux ARM64

```
curl -L https://github.com/hurtki/fed/releases/latest/download/fed-linux-arm64 -o fed
chmod +x fed
sudo mv fed /usr/local/bin/
```

### macOS Apple Silicon(ARM)

```
curl -L https://github.com/hurtki/fed/releases/latest/download/fed-darwin-arm64 -o fed
chmod +x fed
sudo mv fed /usr/local/bin/
```

## Usage

To start the agent, run the `fed` command. You can optionally pass a path to a custom environment file (defaults to `.env`):

```bash
fed [path_to_env_file]
```

Upon startup, the CLI will guide you through the following interactive steps:

1. **LLM Selection**: You will be prompted to choose which LLM provider to use (`gemini` or `ollama`).
2. **Interactive Chat Loop**: Once initialized, the agent enters a continuous loop, prompting you with `->` for input. You can type your instructions or queries (supporting multi-line input). To submit your input, press **Tab** or press **Enter twice in a row**. To exit, simply press Enter on an empty prompt.
