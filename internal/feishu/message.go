package feishu

import (
	"mcp-feishu/internal/types"
)

// MessageBuilder 消息构建器
type MessageBuilder struct {
	securityManager *SecurityManager
}

// NewMessageBuilder 创建消息构建器
func NewMessageBuilder(securityManager *SecurityManager) *MessageBuilder {
	return &MessageBuilder{
		securityManager: securityManager,
	}
}

// BuildTextMessage 构建文本消息
func (mb *MessageBuilder) BuildTextMessage(text string) (*types.FeishuWebhookRequest, error) {
	content := &types.TextMessage{
		Text: text,
	}

	req := &types.FeishuWebhookRequest{
		MsgType: string(types.MessageTypeText),
		Content: content,
	}

	if err := mb.securityManager.ProcessMessage(req, content); err != nil {
		return nil, err
	}

	return req, nil
}

// BuildRichTextMessage 构建富文本消息（Post类型）
func (mb *MessageBuilder) BuildRichTextMessage(content interface{}) (*types.FeishuWebhookRequest, error) {
	postContent := &types.PostMessage{
		Post: content.(map[string]interface{}),
	}

	req := &types.FeishuWebhookRequest{
		MsgType: string(types.MessageTypePost),
		Content: postContent,
	}

	if err := mb.securityManager.ProcessMessage(req, postContent); err != nil {
		return nil, err
	}

	return req, nil
}

// BuildPostMessage 构建结构化富文本消息
func (mb *MessageBuilder) BuildPostMessage(title string, content map[string]interface{}) (*types.FeishuWebhookRequest, error) {
	// 构建完整的post结构
	postData := map[string]interface{}{
		"zh_cn": map[string]interface{}{
			"title":   title,
			"content": content,
		},
	}

	postContent := &types.PostMessage{
		Post: postData,
	}

	req := &types.FeishuWebhookRequest{
		MsgType: string(types.MessageTypePost),
		Content: postContent,
	}

	if err := mb.securityManager.ProcessMessage(req, postContent); err != nil {
		return nil, err
	}

	return req, nil
}

// BuildImageMessage 构建图片消息
func (mb *MessageBuilder) BuildImageMessage(imageKey string) (*types.FeishuWebhookRequest, error) {
	content := &types.ImageMessage{
		ImageKey: imageKey,
	}

	req := &types.FeishuWebhookRequest{
		MsgType: string(types.MessageTypeImage),
		Content: content,
	}

	if err := mb.securityManager.ProcessMessage(req, content); err != nil {
		return nil, err
	}

	return req, nil
}

// BuildInteractiveMessage 构建交互式消息卡片
func (mb *MessageBuilder) BuildInteractiveMessage(config interface{}, elements []interface{}, header interface{}) (*types.FeishuWebhookRequest, error) {
	content := &types.InteractiveMessage{
		Config:   config,
		Elements: elements,
		Header:   header,
	}

	req := &types.FeishuWebhookRequest{
		MsgType: string(types.MessageTypeInteractive),
		Content: content,
	}

	if err := mb.securityManager.ProcessMessage(req, content); err != nil {
		return nil, err
	}

	return req, nil
}

// BuildShareChatMessage 构建群名片消息
func (mb *MessageBuilder) BuildShareChatMessage(shareChatID string) (*types.FeishuWebhookRequest, error) {
	content := &types.ShareChatMessage{
		ShareChatID: shareChatID,
	}

	req := &types.FeishuWebhookRequest{
		MsgType: string(types.MessageTypeShareChat),
		Content: content,
	}

	if err := mb.securityManager.ProcessMessage(req, content); err != nil {
		return nil, err
	}

	return req, nil
}

// CreatePostContent 创建富文本内容的便捷方法
func CreatePostContent(elements ...interface{}) map[string]interface{} {
	return map[string]interface{}{
		"zh_cn": map[string]interface{}{
			"content": [][]interface{}{elements},
		},
	}
}

// CreateTextElement 创建文本元素
func CreateTextElement(text string) map[string]interface{} {
	return map[string]interface{}{
		"tag":  "text",
		"text": text,
	}
}

// CreateLinkElement 创建链接元素
func CreateLinkElement(text, href string) map[string]interface{} {
	return map[string]interface{}{
		"tag":  "a",
		"text": text,
		"href": href,
	}
}

// CreateAtElement 创建@用户元素
func CreateAtElement(userID, userName string) map[string]interface{} {
	return map[string]interface{}{
		"tag":       "at",
		"user_id":   userID,
		"user_name": userName,
	}
}

// CreateImageElement 创建图片元素
func CreateImageElement(imageKey string, width, height int) map[string]interface{} {
	return map[string]interface{}{
		"tag":       "img",
		"image_key": imageKey,
		"width":     width,
		"height":    height,
	}
}

// CreateCardHeader 创建卡片头部
func CreateCardHeader(title, subtitle string, template string) map[string]interface{} {
	return map[string]interface{}{
		"title": map[string]interface{}{
			"tag":     "plain_text",
			"content": title,
		},
		"subtitle": map[string]interface{}{
			"tag":     "plain_text",
			"content": subtitle,
		},
		"template": template,
	}
}

// CreateDivElement 创建分割线元素
func CreateDivElement(text string) map[string]interface{} {
	return map[string]interface{}{
		"tag": "div",
		"text": map[string]interface{}{
			"tag":     "plain_text",
			"content": text,
		},
	}
}

// CreateButtonElement 创建按钮元素
func CreateButtonElement(text, value, actionType string) map[string]interface{} {
	return map[string]interface{}{
		"tag": "button",
		"text": map[string]interface{}{
			"tag":     "plain_text",
			"content": text,
		},
		"value": value,
		"type":  actionType,
	}
}
