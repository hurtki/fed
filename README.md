# Fed - AI Agent

### Future features

- Coding agent ( code writing/refactoring )
- Bash agent ( scripts execution )
- Planning / Thinking

### Ready for now

- Bash agent (scripts execution) is fully implemented and ready for use. Refer to the 'Bash scripts agent' section below for setup and execution instructions.

### Agent Implementation Details

The `Agent` struct is the core component responsible for orchestrating AI-driven tasks. It relies on several interfaces to perform its operations:

- **`AI`**: Provides methods like `GenerateJSON` for interacting with the underlying AI model to generate structured responses from prompts.
- **`ToolChain`**: Offers functionalities such as `RunFileChange` and `ReadFile` to interact with the project's file system, enabling the agent to apply changes and read file content.
- **`AgentReporter`**: Used for reporting the agent's status, logging messages, outlining plans, and communicating the final results of its operations.

#### Bash scripts agent ( golang compiler required )

Setup env variables for gemeni or ollama as in `.env.example` in `.env`

To run agent using gemini:

```
go run ./cmd/gemeni/
```

To run agent using ollama:

```
go run ./cmd/ollama/
```
