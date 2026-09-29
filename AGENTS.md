# AGENTS.md

Guidance and instructions for AI agents working in this repository.

## Project Overview

`ucharla` is a lightweight, zero-dependency CLI REPL written in Go. It connects to local OpenAI-compatible inference servers (e.g., `llama.cpp` server) to provide an interactive terminal chat interface.

## Core Architectural Principles

- **Zero External Dependencies**: Use the Go standard library exclusively. Do not introduce third-party packages without explicit user approval.
- **I/O Separation**: Informational messages, prompts (`> `), and connection statuses write to `os.Stderr`. Generated assistant answers write to `os.Stdout`, allowing responses to be piped cleanly.
- **Stateless by Default**: History is off by default (`-history=false`). When enabled, messages accumulate in memory for the duration of the REPL session.

## Codebase Map

- [`main.go`](main.go): The entire application entry point:
  - `Config`: Runtime configuration struct for server parameters and sampling flags.
  - `parseFlags()`: Flag parsing, default value assignment, and `LLAMA_API_KEY` environment ingestion.
  - `readPrompt()`: Stdin prompt reader until newline or EOF.
  - `ask()`: JSON payload construction, HTTP POST dispatch to `/v1/chat/completions`, error unpacking, and response unmarshaling.
  - `main()`: REPL orchestration loop and history tracking.
- [`Makefile`](Makefile): Build and run automation (`build`, `run`, `clean`, `help`).
- [`go.mod`](go.mod): Go module declaration (`ucharla`).
- [`.editorconfig`](.editorconfig): Formatting rules (tabs for indentation).
- [`.gitattributes`](.gitattributes): Line-ending normalization (`eol=lf`) and Go diff hunk settings.
- [`.gitignore`](.gitignore): Ignored build outputs (`bin/`) and test artifacts.

## Development & Verification Workflow

Always verify your changes before submitting:

1. **Format check**:
   ```bash
   gofmt -d .
   ```
2. **Static analysis**:
   ```bash
   go vet ./...
   ```
3. **Build verification**:
   ```bash
   make build
   ```
4. **Execution check**:
   ```bash
   ./bin/ucharla -h
   ```
5. **Clean up**:
   ```bash
   make clean
   ```

## Commit Guidelines

- Use conventional commits format (e.g., `feat:`, `fix:`, `docs:`, `refactor:`).
- **Never** add "Co-Authored-By" or AI attribution lines to commit messages.
