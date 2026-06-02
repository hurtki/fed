# Fed - AI Agent

> This whole readme is written using this agent

Fed is an AI-driven agent designed to assist with coding, refactoring, and executing bash scripts. It integrates with LLMs (Gemini and Ollama) to understand user prompts, plan tasks, and execute them by interacting with the file system and shell.

## Architecture

The core of the project is the `Agent` struct, which orchestrates tasks using several key interfaces:

- **`AI`**: Interacts with the underlying LLM (Gemini or Ollama) to generate structured JSON responses from prompts.
- **`ToolChain`**: Provides tools for the agent to interact with the environment. It includes functionalities like `RunFileChange` and `ReadFile` to safely read and modify project files, and bash script execution.
- **`AgentReporter`**: Handles reporting the agent's status, logging messages, and communicating the results of its operations to the user.

### Integrations

- **Gemini**: Uses the `google.golang.org/genai` client for interacting with Google's Gemini models.
- **Ollama**: Uses a custom HTTP client to interact with local Ollama instances.

## Internal Workflow

1. **Chat Loop**: The agent runs in a continuous loop, waiting for user input via the CLI.
2. **Planning & Thinking**: Upon receiving a prompt, the agent analyzes the context (project files, chat history) and generates a plan or a direct response.
3. **File Editing**: If the task involves code changes, the agent generates a search/replace patch. The `ToolChain` checks file rights and prompts the user for approval before applying any changes.
4. **Bash Execution**: The agent can propose bash commands to execute. The user is prompted to approve the execution, and the output is fed back into the chat context.

## Getting Started

### Prerequisites

- Go 1.26 or higher
- An API key for Gemini OR a running local instance of Ollama

### Configuration

1. Copy the example environment file:
   ```bash
   cp .env.example .env
   ```
2. Edit `.env` to configure your preferred LLM:
   - For **Ollama**: Set `OLLAMA_ENDPOINT` and `OLLAMA_MODEL`.
   - For **Gemini**: Set `GEMINI_TOKEN` and `GEMINI_MODEL`.

### Building and Running

You can run the agent directly using Go. The agent will prompt you to choose which LLM to use (`gemini` or `ollama`).

**Run the standard agent:**

```bash
go run ./cmd/agent/
```

**Run the development agent (with debug logging):**

```bash
go run ./cmd/agent-dev/
```

Once started, type your prompt at the `->` indicator to interact with the agent.
