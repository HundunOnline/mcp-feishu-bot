package types

// FeishuConfig 飞书配置
type FeishuConfig struct {
	WebhookURL   string   `json:"webhook_url"`
	Secret       string   `json:"secret,omitempty"`   // 签名校验密钥
	Keywords     []string `json:"keywords,omitempty"` // 自定义关键词
	SecurityType string   `json:"security_type"`      // none, signature, keyword
}

// MessageType 消息类型
type MessageType string

const (
	MessageTypeText        MessageType = "text"
	MessageTypePost        MessageType = "post" // 富文本消息
	MessageTypeImage       MessageType = "image"
	MessageTypeInteractive MessageType = "interactive"
	MessageTypeShareChat   MessageType = "share_chat"
)

// SecurityType 安全设置类型
type SecurityType string

const (
	SecurityTypeNone      SecurityType = "none"
	SecurityTypeSignature SecurityType = "signature"
	SecurityTypeKeyword   SecurityType = "keyword"
)

// TextMessage 文本消息
type TextMessage struct {
	Text string `json:"text"`
}

// PostMessage 富文本消息
type PostMessage struct {
	Post map[string]interface{} `json:"post"`
}

// ImageMessage 图片消息
type ImageMessage struct {
	ImageKey string `json:"image_key"`
}

// InteractiveMessage 消息卡片
type InteractiveMessage struct {
	Config   interface{}   `json:"config,omitempty"`
	Elements []interface{} `json:"elements"`
	Header   interface{}   `json:"header,omitempty"`
}

// ShareChatMessage 群名片
type ShareChatMessage struct {
	ShareChatID string `json:"share_chat_id"`
}

// FeishuWebhookRequest 飞书Webhook请求结构
type FeishuWebhookRequest struct {
	MsgType   string      `json:"msg_type"`
	Content   interface{} `json:"content"`
	Timestamp int64       `json:"timestamp,omitempty"`
	Sign      string      `json:"sign,omitempty"`
}

// FeishuWebhookResponse 飞书Webhook响应结构
type FeishuWebhookResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data,omitempty"`
}

// MCPRequest MCP请求结构
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"` // 通知请求可能没有id
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// MCPResponse MCP响应结构
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

// MCPError MCP错误结构
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Tool MCP工具定义
type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// ToolCall 工具调用
type ToolCall struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

// ToolResult 工具结果
type ToolResult struct {
	Content []interface{} `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}
