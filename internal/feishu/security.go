package feishu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"mcp-feishu/internal/types"
	"strconv"
	"strings"
	"time"
)

// SecurityManager 安全管理器
type SecurityManager struct {
	securityType types.SecurityType
	secret       string
	keywords     []string
}

// NewSecurityManager 创建安全管理器
func NewSecurityManager(securityType types.SecurityType, secret string, keywords []string) *SecurityManager {
	return &SecurityManager{
		securityType: securityType,
		secret:       secret,
		keywords:     keywords,
	}
}

// ProcessMessage 处理消息安全设置
func (sm *SecurityManager) ProcessMessage(req *types.FeishuWebhookRequest, content interface{}) error {
	switch sm.securityType {
	case types.SecurityTypeSignature:
		return sm.addSignature(req)
	case types.SecurityTypeKeyword:
		return sm.validateKeyword(content)
	case types.SecurityTypeNone:
		return nil
	default:
		return fmt.Errorf("不支持的安全类型: %s", sm.securityType)
	}
}

// addSignature 添加签名
func (sm *SecurityManager) addSignature(req *types.FeishuWebhookRequest) error {
	if sm.secret == "" {
		return fmt.Errorf("密钥不能为空")
	}

	timestamp := time.Now().Unix()
	req.Timestamp = timestamp

	// 构造签名字符串：timestamp + "\n" + secret
	stringToSign := fmt.Sprintf("%d\n%s", timestamp, sm.secret)

	// 计算HMAC-SHA256签名
	h := hmac.New(sha256.New, []byte(sm.secret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	req.Sign = signature
	return nil
}

// validateKeyword 验证关键词
func (sm *SecurityManager) validateKeyword(content interface{}) error {
	if len(sm.keywords) == 0 {
		return fmt.Errorf("关键词列表不能为空")
	}

	// 提取消息文本内容
	text := extractTextFromContent(content)
	if text == "" {
		return fmt.Errorf("无法提取消息文本内容")
	}

	// 检查是否包含任意一个关键词
	for _, keyword := range sm.keywords {
		if strings.Contains(text, keyword) {
			return nil
		}
	}

	return fmt.Errorf("消息内容必须包含以下关键词之一: %v", sm.keywords)
}

// extractTextFromContent 从不同类型的消息内容中提取文本
func extractTextFromContent(content interface{}) string {
	switch v := content.(type) {
	case *types.TextMessage:
		return v.Text
	case map[string]interface{}:
		// 处理不同消息类型
		if text, ok := v["text"].(string); ok {
			return text
		}
		if title, ok := v["title"].(string); ok {
			return title
		}
		// 尝试递归提取文本
		for _, value := range v {
			if text := extractTextFromContent(value); text != "" {
				return text
			}
		}
	case []interface{}:
		// 处理数组类型
		for _, item := range v {
			if text := extractTextFromContent(item); text != "" {
				return text
			}
		}
	case string:
		return v
	}
	return ""
}

// ValidateSignature 验证接收到的签名（用于接收飞书回调）
func (sm *SecurityManager) ValidateSignature(timestamp string, signature string, body []byte) error {
	if sm.securityType != types.SecurityTypeSignature {
		return nil
	}

	if sm.secret == "" {
		return fmt.Errorf("密钥未配置")
	}

	// 验证时间戳（防重放攻击）
	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return fmt.Errorf("无效的时间戳格式")
	}

	now := time.Now().Unix()
	if now-ts > 300 { // 5分钟有效期
		return fmt.Errorf("请求时间戳过期")
	}

	// 重新计算签名
	stringToSign := fmt.Sprintf("%s\n%s", timestamp, sm.secret)
	h := hmac.New(sha256.New, []byte(stringToSign))
	h.Write(body)
	expectedSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return fmt.Errorf("签名验证失败")
	}

	return nil
}
