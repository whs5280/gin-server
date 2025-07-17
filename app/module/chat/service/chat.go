package service

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"gin-server/app/module/chat/model"
	"io"
	"net/http"
	"os"
	"time"
)

func ValidLength(input string) bool {
	return len(input) >= 3
}

func Reply(input string, method string) string {
	request := buildChatCompletionRequest(input, method)

	if method == model.ChatStream {
		fmt.Println(model.StreamMessage)
		err := chatCompletionStream(request, func(response *model.ChatCompletionResponse) error {
			if len(response.Choices) > 0 {
				fmt.Print(response.Choices[0].Delta.Content) // 注意: 流式响应通常使用 Delta 而非 Message
				time.Sleep(100 * time.Millisecond)           // 模拟流式效果
			}
			return nil
		})
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
		}
	}

	if method == model.NoStream {
		response, err := chatCompletion(request)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return model.ErrorMessage
		}

		if len(response.Choices) > 0 {
			return response.Choices[0].Message.Content
		} else {
			fmt.Printf("\nUsage: Prompt Tokens: %d, Completion Tokens: %d, Total Tokens: %d\n",
				response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens)
			return "No response from API"
		}
	}
	return "end"
}

// buildChatCompletionRequest 构建聊天补全请求
func buildChatCompletionRequest(keyword string, method string) model.ChatCompletionRequest {
	request := model.ChatCompletionRequest{
		Model: model.DefaultModel,
		Messages: []model.Message{
			{
				Role:    "user",
				Content: keyword,
			},
		},
		Temperature: model.Temperature,
		MaxTokens:   model.MaxTokens,
	}

	if method == model.ChatStream {
		request.Stream = true
	}
	return request
}

// chatCompletion 调用 DeepSeek API 进行聊天补全
func chatCompletion(request model.ChatCompletionRequest) (*model.ChatCompletionResponse, error) {
	c := &model.DeepSeekClient{
		ApiKey:     getApiKey(),
		BaseURL:    model.BaseUrl,
		HttpClient: &http.Client{},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %v", err)
	}

	url := c.BaseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %v", err)
	}
	if c.ApiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response model.ChatCompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("decode response failed: %v", err)
	}

	return &response, nil
}

/**
 * bytes.Buffer 是一个实现了 io.Writer 接口的缓冲区，用于存储写入的数据。
 */
func chatCompletionStream(request model.ChatCompletionRequest, callback func(*model.ChatCompletionResponse) error) error {
	c := &model.DeepSeekClient{
		ApiKey:     getApiKey(),
		BaseURL:    model.BaseUrl,
		HttpClient: &http.Client{},
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal request failed: %v", err)
	}

	url := c.BaseURL + "/chat/completions"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return fmt.Errorf("create request failed: %v", err)
	}
	if c.ApiKey == "" {
		return fmt.Errorf("API key is required")
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.ApiKey)
	req.Header.Set("Accept", "text/event-stream") // 重要：声明接受 SSE 格式

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	// 流式处理响应 bufio.Reader 逐行处理 SSE 格式
	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("read response failed: %w", err)
		}

		// 清理 SSE 格式的行
		cleaned := bytes.TrimSpace(line)
		// fmt.Printf("Raw line: %s\n", string(line))  // 调试
		if len(cleaned) == 0 || bytes.HasPrefix(cleaned, []byte(":")) {
			continue // 跳过空行和注释行
		}
		// 处理 data: 前缀
		if bytes.HasPrefix(cleaned, []byte("data: ")) {
			cleaned = bytes.TrimPrefix(cleaned, []byte("data: "))
		}
		// 检查是否是 [DONE] 事件
		if bytes.Equal(cleaned, []byte("[DONE]")) {
			break
		}

		// 解析 JSON
		var response model.ChatCompletionResponse
		if err := json.Unmarshal(cleaned, &response); err != nil {
			fmt.Printf("Failed to parse line: %s\n", string(cleaned))
			return fmt.Errorf("decode response failed: %w, content: %s", err, string(cleaned))
		}

		if err := callback(&response); err != nil {
			return err
		}
	}

	return nil
}

func getApiKey() string {
	return os.Getenv("DEEPSEEK_API_KEY")
}
