package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type Config struct {
	BaseURL         string
	APIKey          string
	Model           string
	Temperature     float64
	TopP            float64
	TopK            int
	MinP            float64
	MaxTokens       int
	Seed            int
	ReasoningEffort string
	Timeout         time.Duration
	KeepHistory     bool
	System          string
}

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model           string        `json:"model,omitempty"`
	Messages        []ChatMessage `json:"messages"`
	Temperature     float64       `json:"temperature"`
	TopP            float64       `json:"top_p"`
	TopK            int           `json:"top_k,omitempty"`
	MinP            float64       `json:"min_p,omitempty"`
	MaxTokens       int           `json:"max_tokens"`
	Seed            int           `json:"seed,omitempty"`
	ReasoningEffort string        `json:"reasoning_effort,omitempty"`
	Stream          bool          `json:"stream"`
}

type ChatResponse struct {
	Choices []struct {
		Index   int         `json:"index"`
		Message ChatMessage `json:"message"`
	} `json:"choices"`
}

type APIError struct {
	Error struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    any    `json:"code"`
	} `json:"error"`
}

const (
	defaultBaseURL         = "http://localhost:8080"
	defaultTemperature     = 0.1
	defaultTopP            = 0.9
	defaultTopK            = 40
	defaultMinP            = 0.05
	defaultMaxTokens       = 1200
	defaultSeed            = 42
	defaultTimeout         = 10 * time.Minute
	defaultReasoningEffort = "none"

	apiKeyEnv = "LLAMA_API_KEY"
)

func main() {
	cfg := parseFlags()

	client := &http.Client{
		Timeout: cfg.Timeout,
	}

	messages := make([]ChatMessage, 0, 16)

	if cfg.System != "" {
		messages = append(messages, ChatMessage{
			Role:    "system",
			Content: cfg.System,
		})
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Fprintln(os.Stderr, "llama.cpp REPL")
	fmt.Fprintf(
		os.Stderr,
		"URL: %s/v1/chat/completions\n",
		strings.TrimRight(cfg.BaseURL, "/"),
	)

	fmt.Fprintf(
		os.Stderr,
		"temperature=%g top_p=%g top_k=%d min_p=%g max_tokens=%d seed=%d history=%t\n",
		cfg.Temperature,
		cfg.TopP,
		cfg.TopK,
		cfg.MinP,
		cfg.MaxTokens,
		cfg.Seed,
		cfg.KeepHistory,
	)

	if cfg.APIKey != "" {
		fmt.Fprintf(os.Stderr, "API key: loaded from %s\n", apiKeyEnv)
	} else {
		fmt.Fprintf(os.Stderr, "API key: %s is not set\n", apiKeyEnv)
	}

	fmt.Fprintln(
		os.Stderr,
		"Enter a prompt. EOF exits (^D on Unix, ^Z then Enter on Windows).",
	)
	fmt.Fprintln(os.Stderr)

	for {
		fmt.Fprint(os.Stderr, "> ")

		prompt, err := readPrompt(reader)
		if err != nil && !errors.Is(err, io.EOF) {
			fmt.Fprintf(os.Stderr, "input error: %v\n", err)
			os.Exit(1)
		}

		prompt = strings.TrimSpace(prompt)

		if prompt != "" {
			requestMessages := messages

			if !cfg.KeepHistory {
				requestMessages = nil

				if cfg.System != "" {
					requestMessages = append(
						requestMessages,
						ChatMessage{
							Role:    "system",
							Content: cfg.System,
						},
					)
				}
			}

			requestMessages = append(
				requestMessages,
				ChatMessage{
					Role:    "user",
					Content: prompt,
				},
			)

			answer, requestErr := ask(
				client,
				cfg,
				requestMessages,
			)

			if requestErr != nil {
				fmt.Fprintf(
					os.Stderr,
					"request error: %v\n",
					requestErr,
				)
			} else {
				fmt.Println(answer)
				fmt.Println()

				if cfg.KeepHistory {
					messages = requestMessages

					messages = append(
						messages,
						ChatMessage{
							Role:    "assistant",
							Content: answer,
						},
					)
				}
			}
		}

		if errors.Is(err, io.EOF) {
			fmt.Fprintln(os.Stderr)
			return
		}
	}
}

func parseFlags() Config {
	var cfg Config

	flag.StringVar(
		&cfg.BaseURL,
		"url",
		defaultBaseURL,
		"llama.cpp base URL",
	)

	flag.StringVar(
		&cfg.Model,
		"model",
		"",
		"model name sent in the request",
	)

	flag.Float64Var(
		&cfg.Temperature,
		"temperature",
		defaultTemperature,
		"sampling temperature",
	)

	flag.Float64Var(
		&cfg.TopP,
		"top-p",
		defaultTopP,
		"top-p sampling value",
	)

	flag.IntVar(
		&cfg.TopK,
		"top-k",
		defaultTopK,
		"top-k sampling value",
	)

	flag.Float64Var(
		&cfg.MinP,
		"min-p",
		defaultMinP,
		"min-p sampling value",
	)

	flag.IntVar(
		&cfg.MaxTokens,
		"max-tokens",
		defaultMaxTokens,
		"maximum number of generated tokens",
	)

	flag.IntVar(
		&cfg.Seed,
		"seed",
		defaultSeed,
		"sampling seed",
	)

	flag.StringVar(
		&cfg.ReasoningEffort,
		"reasoning-effort",
		defaultReasoningEffort,
		"Model reasoning effort label",
	)

	flag.DurationVar(
		&cfg.Timeout,
		"timeout",
		defaultTimeout,
		"HTTP request timeout",
	)

	flag.BoolVar(
		&cfg.KeepHistory,
		"history",
		false,
		"keep previous user/assistant messages in the conversation",
	)

	flag.StringVar(
		&cfg.System,
		"system",
		"",
		"optional system prompt",
	)

	flag.Parse()

	cfg.APIKey = strings.TrimSpace(os.Getenv(apiKeyEnv))

	return cfg
}

func readPrompt(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	return line, err
}

func ask(
	client *http.Client,
	cfg Config,
	messages []ChatMessage,
) (string, error) {
	reqBody := ChatRequest{
		Model:           cfg.Model,
		Messages:        messages,
		Temperature:     cfg.Temperature,
		TopP:            cfg.TopP,
		TopK:            cfg.TopK,
		MinP:            cfg.MinP,
		MaxTokens:       cfg.MaxTokens,
		Seed:            cfg.Seed,
		Stream:          false,
		ReasoningEffort: cfg.ReasoningEffort,
	}

	data, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf(
			"marshal request: %w",
			err,
		)
	}

	endpoint :=
		strings.TrimRight(cfg.BaseURL, "/") +
			"/v1/chat/completions"

	req, err := http.NewRequest(
		http.MethodPost,
		endpoint,
		bytes.NewReader(data),
	)
	if err != nil {
		return "", fmt.Errorf(
			"create HTTP request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	if cfg.APIKey != "" {
		req.Header.Set(
			"Authorization",
			"Bearer "+cfg.APIKey,
		)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf(
			"send HTTP request: %w",
			err,
		)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf(
			"read response: %w",
			err,
		)
	}

	if resp.StatusCode < 200 ||
		resp.StatusCode >= 300 {

		var apiErr APIError

		if json.Unmarshal(body, &apiErr) == nil &&
			apiErr.Error.Message != "" {

			return "", fmt.Errorf(
				"HTTP %s: %s",
				resp.Status,
				apiErr.Error.Message,
			)
		}

		return "", fmt.Errorf(
			"HTTP %s: %s",
			resp.Status,
			strings.TrimSpace(string(body)),
		)
	}

	var chatResp ChatResponse

	if err := json.Unmarshal(
		body,
		&chatResp,
	); err != nil {
		return "", fmt.Errorf(
			"decode response: %w\nraw response: %s",
			err,
			string(body),
		)
	}

	if len(chatResp.Choices) == 0 {
		return "", errors.New(
			"server returned no choices",
		)
	}

	return chatResp.Choices[0].Message.Content, nil
}
