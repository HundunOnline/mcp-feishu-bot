package feishu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mcp-feishu/internal/types"
	"net/http"
	"time"
)

// Client 飞书客户端
type Client struct {
	webhookURL      string
	httpClient      *http.Client
	messageBuilder  *MessageBuilder
	securityManager *SecurityManager
}

// NewClient 创建飞书客户端
func NewClient(config types.FeishuConfig) *Client {
	securityManager := NewSecurityManager(
		types.SecurityType(config.SecurityType),
		config.Secret,
		config.Keywords,
	)

	return &Client{
		webhookURL: config.WebhookURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		messageBuilder:  NewMessageBuilder(securityManager),
		securityManager: securityManager,
	}
}

// SendMessage 发送消息
func (c *Client) SendMessage(req *types.FeishuWebhookRequest) (*types.FeishuWebhookResponse, error) {
	// 序列化请求
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequest("POST", c.webhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建HTTP请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送HTTP请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var feishuResp types.FeishuWebhookResponse
	if err := json.Unmarshal(body, &feishuResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, 响应内容: %s", err, string(body))
	}

	// 检查响应状态
	if feishuResp.Code != 0 {
		return &feishuResp, fmt.Errorf("飞书API返回错误: code=%d, message=%s", feishuResp.Code, feishuResp.Message)
	}

	return &feishuResp, nil
}

// SendTextMessage 发送文本消息
func (c *Client) SendTextMessage(text string) (*types.FeishuWebhookResponse, error) {
	req, err := c.messageBuilder.BuildTextMessage(text)
	if err != nil {
		return nil, fmt.Errorf("构建文本消息失败: %w", err)
	}

	return c.SendMessage(req)
}

// SendRichTextMessage 发送富文本消息
func (c *Client) SendRichTextMessage(content interface{}) (*types.FeishuWebhookResponse, error) {
	req, err := c.messageBuilder.BuildRichTextMessage(content)
	if err != nil {
		return nil, fmt.Errorf("构建富文本消息失败: %w", err)
	}

	return c.SendMessage(req)
}

// SendPostMessage 发送群名片消息
func (c *Client) SendPostMessage(title string, content map[string]interface{}) (*types.FeishuWebhookResponse, error) {
	req, err := c.messageBuilder.BuildPostMessage(title, content)
	if err != nil {
		return nil, fmt.Errorf("构建群名片消息失败: %w", err)
	}

	return c.SendMessage(req)
}

// SendImageMessage 发送图片消息
func (c *Client) SendImageMessage(imageKey string) (*types.FeishuWebhookResponse, error) {
	req, err := c.messageBuilder.BuildImageMessage(imageKey)
	if err != nil {
		return nil, fmt.Errorf("构建图片消息失败: %w", err)
	}

	return c.SendMessage(req)
}

// SendInteractiveMessage 发送交互式消息卡片
func (c *Client) SendInteractiveMessage(config interface{}, elements []interface{}, header interface{}) (*types.FeishuWebhookResponse, error) {
	req, err := c.messageBuilder.BuildInteractiveMessage(config, elements, header)
	if err != nil {
		return nil, fmt.Errorf("构建交互式消息失败: %w", err)
	}

	return c.SendMessage(req)
}

// SendShareChatMessage 发送群名片消息
func (c *Client) SendShareChatMessage(shareChatID string) (*types.FeishuWebhookResponse, error) {
	req, err := c.messageBuilder.BuildShareChatMessage(shareChatID)
	if err != nil {
		return nil, fmt.Errorf("构建群名片消息失败: %w", err)
	}

	return c.SendMessage(req)
}

// GetSecurityManager 获取安全管理器
func (c *Client) GetSecurityManager() *SecurityManager {
	return c.securityManager
}

// UpdateConfig 更新配置
func (c *Client) UpdateConfig(config types.FeishuConfig) {
	c.webhookURL = config.WebhookURL
	c.securityManager = NewSecurityManager(
		types.SecurityType(config.SecurityType),
		config.Secret,
		config.Keywords,
	)
	c.messageBuilder = NewMessageBuilder(c.securityManager)
}
