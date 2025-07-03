package model

import "time"

const (
	BaseDelay = 3 * time.Second
	FilePath  = "data/training.jsonl"
)

type OpenAIFile struct {
	ID            string `json:"id"`
	Object        string `json:"object"`
	Bytes         int    `json:"bytes"`
	CreatedAt     int64  `json:"created_at"`
	Filename      string `json:"filename"`
	Purpose       string `json:"purpose"`
	Status        string `json:"status,omitempty"`
	StatusDetails string `json:"status_details,omitempty"`
}

type FineTuningJobResponse struct {
	ID             string `json:"id"` // 对象标识符,可以在API端点中引用
	Object         string `json:"object"`
	Model          string `json:"model"` // 被微调的基础模型
	CreatedAt      int64  `json:"created_at"`
	Status         string `json:"status"`           // 微调作业的当前状态,可以是validating_files、queued、running、succeeded、failed或cancelled
	FineTunedModel string `json:"fine_tuned_model"` // 正在创建的微调模型的名称。如果微调作业仍在运行,则值为null
}
