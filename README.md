# ucharla

A lightweight, zero-dependency Go CLI REPL for interacting with local OpenAI-compatible LLM servers (such as `llama.cpp` server, LocalAI, or Ollama) via the `/v1/chat/completions` endpoint.

## Features

- **Zero External Dependencies**: Built entirely with the Go standard library.
- **Interactive REPL**: Read-eval-print loop with EOF support (`Ctrl+D` on Unix, `Ctrl+Z` then Enter on Windows).
- **Optional Conversation History**: Toggle multi-turn memory (`-history`) or treat each query as stateless.
- **Granular Sampling Control**: Full control over temperature, top-p, top-k, min-p, max tokens, seed, and reasoning effort.
- **Authentication Support**: Reads `LLAMA_API_KEY` from the environment when authentication is enabled on the server.
- **Makefile Automation**: Clean, standardized build and run targets.

## Prerequisites

- [Go](https://go.dev/) (1.21+ recommended)
- A running OpenAI-compatible server (e.g., `llama-server` listening on `http://localhost:8080`)

## Getting Started

### Build

Compile the binary to `bin/ucharla`:

```bash
make build
```

### Run

Launch the interactive REPL:

```bash
make run
```

Pass custom flags to the application:

```bash
make run ARGS="-url http://localhost:8080 -model llama-3 -temperature 0.7 -history"
```

Or execute the compiled binary directly:

```bash
./bin/ucharla -url http://localhost:8080 -history
```

### Code Quality

Format code and run static analysis:

```bash
make fmt   # Formats Go code using gofmt -s -w
make vet   # Runs go vet ./...
make check # Runs both fmt and vet
```

### Clean

Remove build artifacts:

```bash
make clean
```

## Configuration Flags

| Flag | Type | Default | Description |
| :--- | :--- | :--- | :--- |
| `-url` | string | `http://localhost:8080` | Base URL of the OpenAI-compatible server |
| `-model` | string | `""` | Model identifier sent in the request |
| `-system` | string | `""` | Optional system prompt to prepend |
| `-history` | bool | `false` | Retain conversation history across turns |
| `-temperature` | float | `0.1` | Sampling temperature |
| `-top-p` | float | `0.9` | Top-p (nucleus) sampling cutoff |
| `-top-k` | int | `40` | Top-k sampling limit |
| `-min-p` | float | `0.05` | Min-p sampling threshold |
| `-max-tokens` | int | `1200` | Maximum generated tokens per response |
| `-seed` | int | `42` | Random seed for deterministic generation |
| `-reasoning-effort` | string | `"none"` | Reasoning effort label (e.g. `none`, `low`, `medium`, `high`) |
| `-timeout` | duration | `10m0s` | HTTP client request timeout |

### Environment Variables

- `LLAMA_API_KEY`: If set, included as a Bearer token in the `Authorization` header (`Authorization: Bearer <key>`).
