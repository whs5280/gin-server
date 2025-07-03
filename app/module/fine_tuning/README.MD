- API文档： `https://openai.apifox.cn/doc-3222700`

## 1. `purpose: "fine-tune"`（用于模型微调）
### 用途
- 基于 `gpt-3.5-turbo` 训练专有领域模型（客服/法律/医疗等）
- 调整模型行为使其符合特定需求

### 文件格式要求
- **必须使用 JSONL 格式**（每行一个完整JSON对象）
- **必须包含** `messages` 消息结构

#### 示例训练数据
```jsonl
{"messages": [{"role": "system", "content": "你是一个客服助手"}, {"role": "user", "content": "如何退款？"}, {"role": "assistant", "content": "请提供订单号。"}]}
{"messages": [{"role": "system", "content": "你是一个客服助手"}, {"role": "user", "content": "我的订单没收到"}, {"role": "assistant", "content": "请检查物流信息。"}]}
```

## 2. `purpose: "assistants"`（用于助手和消息）
### 用途
- 上传知识库文件（`PDF`、`TXT`），让助手能基于文件内容回答问题
- 动态扩展助手的知识，但不修改模型本身

### 文件格式要求
- **PDF/TXT/DOCX**