package model

import "net/http"

const (
	EXIT           = "退出"
	WelcomeMessage = "尽管问..."
	ExitMessage    = "Bye！"
	ValidMessage   = "长度不能小于3"
	ReplyMessage   = "AI: %s\n"
	ErrorMessage   = "Api 请求失败"
	ChatStream     = "stream"
	NoStream       = "normal"
	StreamMessage  = "开始接收流式响应:"
	BaseUrl        = "https://api.gpt.ge/v1"
	DefaultModel   = "gpt-4o"
	Temperature    = 0.7
	MaxTokens      = 200
)

type DeepSeekClient struct {
	ApiKey     string
	BaseURL    string
	HttpClient *http.Client
}

// ChatCompletionRequest 聊天补全请求结构体
type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse 聊天补全响应结构体
type ChatCompletionResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	// Choices 是一个数组，包含了多个选项
	Choices []struct {
		Index        int `json:"index"`
		Message      Message
		Delta        Message
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	// Usage 记录了本次对话的使用情况
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
}
