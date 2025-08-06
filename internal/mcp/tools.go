package mcp

import (
	"encoding/json"
	"fmt"
	"mcp-feishu/internal/feishu"
	"mcp-feishu/internal/types"
)

// ToolsHandler 工具处理器
type ToolsHandler struct {
	feishuClient *feishu.Client
}

// NewToolsHandler 创建工具处理器
func NewToolsHandler(feishuClient *feishu.Client) *ToolsHandler {
	return &ToolsHandler{
		feishuClient: feishuClient,
	}
}

// GetTools 获取所有可用工具
func (th *ToolsHandler) GetTools() []types.Tool {
	return []types.Tool{
		{
			Name:        "send_text_message",
			Description: "发送纯文本消息\n\n发送简单的文本消息，支持自动换行和基本文本格式。这是最基础的消息类型，适合发送通知、状态更新等简单信息。会自动应用配置的安全设置。\n\n示例：{\"text\": \"系统维护通知：服务将在今晚22:00-24:00进行维护\"}",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"text": map[string]interface{}{
						"type":        "string",
						"description": "要发送的纯文本内容，支持换行符。最大长度为30000字符。如果配置了关键词验证，文本必须包含指定关键词。",
					},
				},
				"required": []string{"text"},
			},
		},
		{
			Name:        "send_post_message",
			Description: "发送富文本消息\n\n发送支持格式化的富文本消息，可以包含加粗、斜体、链接、@用户、图片等丰富格式。支持可选标题。AI只需提供内容数组和可选标题，工具会自动包装成飞书API格式。\n\n示例1（无标题）：{\"content\": [[{\"tag\": \"text\", \"text\": \"访问\"}, {\"tag\": \"a\", \"text\": \"飞书官网\", \"href\": \"https://feishu.cn\"}]]}\n示例2（有标题）：{\"title\": \"数据报告\", \"content\": [[{\"tag\": \"text\", \"text\": \"本月销售额增长15%\"}]]}",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"title": map[string]interface{}{
						"type":        "string",
						"description": "可选的消息标题，会显示在消息顶部。如果不提供，则发送无标题的富文本消息。建议不超过100字符。",
					},
					"content": map[string]interface{}{
						"type":        "array",
						"description": "富文本内容数组，二维数组格式。每行是一个元素数组，元素可以是文本、链接、@用户等。支持的元素类型：text(文本)、a(链接)、at(提及用户)、img(图片)等。工具会自动包装成飞书API需要的post结构。",
					},
				},
				"required": []string{"content"},
			},
		},
		{
			Name:        "send_image_message",
			Description: "发送图片消息\n\n发送图片到飞书群组或个人。需要先通过飞书API上传图片获取image_key，然后使用此工具发送。支持常见图片格式。\n\n注意：image_key必须是有效的飞书图片资源标识符。\n\n示例：{\"image_key\": \"img_v2_041b28e3-xxxx-xxxx-xxxx-xxxxxxxxxxxx\"}",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"image_key": map[string]interface{}{
						"type":        "string",
						"description": "飞书图片资源的唯一标识符，格式通常为 img_v2_ 开头的字符串。需要先通过飞书上传图片接口获取此值。image_key有效期通常为24小时。",
					},
				},
				"required": []string{"image_key"},
			},
		},
		{
			Name:        "send_interactive_message",
			Description: "发送交互式消息卡片\n\n发送功能丰富的交互式卡片，可以包含按钮、表单、布局等复杂UI元素。适合发送需要用户交互的通知、审批、问卷等场景。\n\n示例：{\"elements\": [{\"tag\": \"div\", \"text\": {\"tag\": \"plain_text\", \"content\": \"请确认操作\"}}, {\"tag\": \"button\", \"text\": {\"tag\": \"plain_text\", \"content\": \"确认\"}, \"type\": \"primary\"}]}",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"config": map[string]interface{}{
						"type":        "object",
						"description": "可选的卡片全局配置，如宽度模式、更新设置等。格式：{\"wide_screen_mode\": true, \"enable_forward\": false}",
					},
					"elements": map[string]interface{}{
						"type":        "array",
						"description": "卡片主要内容元素数组，每个元素代表卡片中的一个组件。支持的元素类型：div(文本块)、hr(分割线)、img(图片)、action(按钮组)、field(字段组)等。每个元素必须包含tag字段指定类型。",
					},
					"header": map[string]interface{}{
						"type":        "object",
						"description": "可选的卡片头部配置，包含标题、副标题、模板样式等。格式：{\"title\": {\"tag\": \"plain_text\", \"content\": \"标题\"}, \"template\": \"blue\"}。template可选值：blue、wathet、turquoise、green、yellow、orange、red、carmine、violet、purple、indigo、grey",
					},
				},
				"required": []string{"elements"},
			},
		},
		{
			Name:        "send_share_chat_message",
			Description: "发送群聊分享卡片\n\n分享一个群聊的卡片信息到当前会话，类似于微信的群聊名片功能。接收者可以通过卡片快速加入指定群组。\n\n注意：share_chat_id必须是有效的飞书群组ID，且当前机器人必须有权限访问该群组。\n\n示例：{\"share_chat_id\": \"oc_a0553eda9014c201e6969b478895c230\"}",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"share_chat_id": map[string]interface{}{
						"type":        "string",
						"description": "要分享的群聊的唯一标识符，格式通常为 oc_ 开头的字符串。可以通过飞书群聊设置或API获取。机器人必须是该群聊的成员才能分享。",
					},
				},
				"required": []string{"share_chat_id"},
			},
		},
	}
}

// CallTool 调用工具
func (th *ToolsHandler) CallTool(toolCall types.ToolCall) (types.ToolResult, error) {
	switch toolCall.Name {
	case "send_text_message":
		return th.handleSendTextMessage(toolCall.Arguments)
	case "send_post_message":
		return th.handleSendPostMessage(toolCall.Arguments)
	case "send_image_message":
		return th.handleSendImageMessage(toolCall.Arguments)
	case "send_interactive_message":
		return th.handleSendInteractiveMessage(toolCall.Arguments)
	case "send_share_chat_message":
		return th.handleSendShareChatMessage(toolCall.Arguments)
	default:
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("未知工具: %s", toolCall.Name),
				},
			},
		}, nil
	}
}

// handleSendTextMessage 处理发送文本消息
func (th *ToolsHandler) handleSendTextMessage(args map[string]interface{}) (types.ToolResult, error) {
	text, ok := args["text"].(string)
	if !ok {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "text 参数必须是字符串类型",
				},
			},
		}, nil
	}

	resp, err := th.feishuClient.SendTextMessage(text)
	if err != nil {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("发送文本消息失败: %v", err),
				},
			},
		}, nil
	}

	return types.ToolResult{
		Content: []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("文本消息发送成功! 响应: code=%d, message=%s", resp.Code, resp.Message),
			},
		},
	}, nil
}

// handleSendPostMessage 处理发送富文本消息
// AI只需提供内容数组和可选标题，工具内部自动包装成完整结构
func (th *ToolsHandler) handleSendPostMessage(args map[string]interface{}) (types.ToolResult, error) {
	// title是可选参数
	title, _ := args["title"].(string)

	content, ok := args["content"]
	if !ok {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "content 参数是必需的",
				},
			},
		}, nil
	}

	// 内部自动包装成完整的post结构
	var postData map[string]interface{}
	if title != "" {
		// 有标题的富文本消息
		postData = map[string]interface{}{
			"zh_cn": map[string]interface{}{
				"title":   title,
				"content": content, // AI直接提供的内容数组
			},
		}
	} else {
		// 无标题的富文本消息
		postData = map[string]interface{}{
			"zh_cn": map[string]interface{}{
				"content": content, // AI直接提供的内容数组
			},
		}
	}

	resp, err := th.feishuClient.SendRichTextMessage(postData)
	if err != nil {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("发送富文本消息失败: %v", err),
				},
			},
		}, nil
	}

	return types.ToolResult{
		Content: []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("富文本消息发送成功! 响应: code=%d, message=%s", resp.Code, resp.Message),
			},
		},
	}, nil
}

// handleSendImageMessage 处理发送图片消息
func (th *ToolsHandler) handleSendImageMessage(args map[string]interface{}) (types.ToolResult, error) {
	imageKey, ok := args["image_key"].(string)
	if !ok {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "image_key 参数必须是字符串类型",
				},
			},
		}, nil
	}

	resp, err := th.feishuClient.SendImageMessage(imageKey)
	if err != nil {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("发送图片消息失败: %v", err),
				},
			},
		}, nil
	}

	return types.ToolResult{
		Content: []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("图片消息发送成功! 响应: code=%d, message=%s", resp.Code, resp.Message),
			},
		},
	}, nil
}

// handleSendInteractiveMessage 处理发送交互式消息
func (th *ToolsHandler) handleSendInteractiveMessage(args map[string]interface{}) (types.ToolResult, error) {
	elements, ok := args["elements"].([]interface{})
	if !ok {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "elements 参数必须是数组类型",
				},
			},
		}, nil
	}

	config := args["config"]
	header := args["header"]

	resp, err := th.feishuClient.SendInteractiveMessage(config, elements, header)
	if err != nil {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("发送交互式消息失败: %v", err),
				},
			},
		}, nil
	}

	return types.ToolResult{
		Content: []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("交互式消息发送成功! 响应: code=%d, message=%s", resp.Code, resp.Message),
			},
		},
	}, nil
}

// handleSendShareChatMessage 处理发送群名片消息
func (th *ToolsHandler) handleSendShareChatMessage(args map[string]interface{}) (types.ToolResult, error) {
	shareChatID, ok := args["share_chat_id"].(string)
	if !ok {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": "share_chat_id 参数必须是字符串类型",
				},
			},
		}, nil
	}

	resp, err := th.feishuClient.SendShareChatMessage(shareChatID)
	if err != nil {
		return types.ToolResult{
			IsError: true,
			Content: []interface{}{
				map[string]interface{}{
					"type": "text",
					"text": fmt.Sprintf("发送群名片消息失败: %v", err),
				},
			},
		}, nil
	}

	return types.ToolResult{
		Content: []interface{}{
			map[string]interface{}{
				"type": "text",
				"text": fmt.Sprintf("群名片消息发送成功! 响应: code=%d, message=%s", resp.Code, resp.Message),
			},
		},
	}, nil
}

// CreateCardElements 创建简单的卡片元素
func CreateCardElements(title, content string) []interface{} {
	return []interface{}{
		feishu.CreateDivElement(title),
		feishu.CreateDivElement(content),
	}
}

// SerializeForLogging 序列化对象用于日志记录
func SerializeForLogging(obj interface{}) string {
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return fmt.Sprintf("序列化失败: %v", err)
	}
	return string(data)
}
