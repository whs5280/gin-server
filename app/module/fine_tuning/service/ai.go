package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gin-server/app/module/fine_tuning/model"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

const (
	customerModelName = "ydt-customer"
	tuningJobsURL     = "https://api.openai.com/v1/fine_tuning/jobs"
	tuningFilesURL    = "https://api.openai.com/v1/files"
	defaultModelName  = "gpt-3.5-turbo"
)

// UploadFile 单个文件的大小最大为 512MB; Fine-tuning API 仅支持jsonl文件
func UploadFile(filePath string, purpose string) (*model.OpenAIFile, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	part, err := writer.CreateFormFile("file", filepath.Base(file.Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file content: %v", err)
	}
	// “fine-tune”进行微调; “assistants”进行助手和消息; 区别见 README.md
	_ = writer.WriteField("purpose", purpose)

	// 关闭writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", tuningFilesURL, payload)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// 处理响应
	if resp.StatusCode != http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(responseBody))
	}

	// 解析响应
	var openAIFile model.OpenAIFile
	if err := json.NewDecoder(resp.Body).Decode(&openAIFile); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &openAIFile, nil
}

// CreatedJob 创建fine-tune任务
func CreatedJob(trainingFileId string) (*model.FineTuningJobResponse, error) {
	// body 参数
	requestData := map[string]interface{}{
		"training_file": trainingFileId,
		"model":         defaultModelName,
		"suffix":        customerModelName,
	}
	jsonData, err := json.Marshal(requestData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request data: %v", err)
	}

	req, err := http.NewRequest("POST", tuningJobsURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getApiKey())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	// 响应体
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var jobResponse model.FineTuningJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	log.Printf("Successfully created fine-tuning job: %s (ID: %s)", jobResponse.Model, jobResponse.ID)
	return &jobResponse, nil
}

// GetJob 检索fine-tune任务详情
func GetJob(jobId string) (*model.FineTuningJobResponse, error) {
	url := fmt.Sprintf("%s/%s", tuningJobsURL, jobId)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+getApiKey())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var jobResponse model.FineTuningJobResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	log.Printf("fine-tuning job status: %s (TunedModel: %s)", jobResponse.Status, jobResponse.FineTunedModel)
	return &jobResponse, nil
}

// 获取API密钥
func getApiKey() string {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is not set")
	}
	return apiKey
}
