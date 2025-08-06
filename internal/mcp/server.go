package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"mcp-feishu/internal/feishu"
	"mcp-feishu/internal/types"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Server MCP服务器
type Server struct {
	feishuClient *feishu.Client
	toolsHandler *ToolsHandler
	logger       zerolog.Logger
}

// NewServer 创建MCP服务器
func NewServer(feishuClient *feishu.Client) *Server {
	// 配置日志
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	toolsHandler := NewToolsHandler(feishuClient)

	return &Server{
		feishuClient: feishuClient,
		toolsHandler: toolsHandler,
		logger:       log.With().Str("component", "mcp-server").Logger(),
	}
}

// Run 运行MCP服务器
func (s *Server) Run() error {
	s.logger.Info().Msg("启动MCP飞书服务器")

	// 设置输入输出
	input := os.Stdin
	output := os.Stdout

	// 处理请求循环
	decoder := json.NewDecoder(input)
	encoder := json.NewEncoder(output)

	for {
		var request types.MCPRequest
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				s.logger.Info().Msg("收到EOF，退出服务器")
				break
			}
			s.logger.Error().Err(err).Msg("解析请求失败")
			continue
		}

		s.logger.Debug().
			Str("method", request.Method).
			Interface("id", request.ID).
			Msg("收到MCP请求")

		response := s.handleRequest(request)

		// 只有非通知请求才需要发送响应
		if response != nil {
			if err := encoder.Encode(response); err != nil {
				s.logger.Error().Err(err).Msg("编码响应失败")
				continue
			}

			s.logger.Debug().
				Interface("id", response.ID).
				Bool("has_error", response.Error != nil).
				Msg("发送MCP响应")
		} else {
			s.logger.Debug().
				Str("method", request.Method).
				Msg("处理通知完成，无响应")
		}
	}

	return nil
}

// handleRequest 处理MCP请求
func (s *Server) handleRequest(request types.MCPRequest) *types.MCPResponse {
	switch request.Method {
	case "initialize":
		response := s.handleInitialize(request)
		return &response
	case "notifications/initialized":
		s.handleInitialized(request)
		return nil // 通知不需要响应
	case "tools/list":
		response := s.handleToolsList(request)
		return &response
	case "tools/call":
		response := s.handleToolsCall(request)
		return &response
	case "ping":
		response := s.handlePing(request)
		return &response
	default:
		response := types.MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &types.MCPError{
				Code:    -32601,
				Message: fmt.Sprintf("方法未找到: %s", request.Method),
			},
		}
		return &response
	}
}

// handleInitialize 处理初始化请求
func (s *Server) handleInitialize(request types.MCPRequest) types.MCPResponse {
	s.logger.Info().Msg("处理初始化请求")

	result := map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities": map[string]interface{}{
			"tools": map[string]interface{}{
				"listChanged": false,
			},
		},
		"serverInfo": map[string]interface{}{
			"name":    "mcp-feishu",
			"version": "1.0.0",
		},
		"instructions": "这是一个飞书消息发送MCP服务器，支持发送各种类型的飞书消息，包括文本、富文本、图片、卡片等。支持三种安全设置：无安全、签名校验、自定义关键词。",
	}

	return types.MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

// handleInitialized 处理初始化完成通知（通知不需要响应）
func (s *Server) handleInitialized(request types.MCPRequest) {
	s.logger.Info().Msg("收到初始化完成通知")
	// 通知处理完成，无需返回响应
}

// handleToolsList 处理工具列表请求
func (s *Server) handleToolsList(request types.MCPRequest) types.MCPResponse {
	s.logger.Info().Msg("处理工具列表请求")

	tools := s.toolsHandler.GetTools()

	result := map[string]interface{}{
		"tools": tools,
	}

	return types.MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

// handleToolsCall 处理工具调用请求
func (s *Server) handleToolsCall(request types.MCPRequest) types.MCPResponse {
	s.logger.Info().Msg("处理工具调用请求")

	// 解析参数
	paramsBytes, err := json.Marshal(request.Params)
	if err != nil {
		return types.MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "参数解析失败",
				Data:    err.Error(),
			},
		}
	}

	var params struct {
		Name      string                 `json:"name"`
		Arguments map[string]interface{} `json:"arguments"`
	}

	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		return types.MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &types.MCPError{
				Code:    -32602,
				Message: "参数格式无效",
				Data:    err.Error(),
			},
		}
	}

	s.logger.Info().
		Str("tool_name", params.Name).
		Interface("arguments", params.Arguments).
		Msg("调用工具")

	// 调用工具
	toolCall := types.ToolCall{
		Name:      params.Name,
		Arguments: params.Arguments,
	}

	result, err := s.toolsHandler.CallTool(toolCall)
	if err != nil {
		return types.MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Error: &types.MCPError{
				Code:    -32603,
				Message: "工具调用失败",
				Data:    err.Error(),
			},
		}
	}

	return types.MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  result,
	}
}

// handlePing 处理ping请求
func (s *Server) handlePing(request types.MCPRequest) types.MCPResponse {
	s.logger.Debug().Msg("处理ping请求")

	return types.MCPResponse{
		JSONRPC: "2.0",
		ID:      request.ID,
		Result:  map[string]interface{}{"status": "pong"},
	}
}

// UpdateFeishuClient 更新飞书客户端
func (s *Server) UpdateFeishuClient(feishuClient *feishu.Client) {
	s.feishuClient = feishuClient
	s.toolsHandler = NewToolsHandler(feishuClient)
	s.logger.Info().Msg("飞书客户端配置已更新")
}

// GetFeishuClient 获取飞书客户端
func (s *Server) GetFeishuClient() *feishu.Client {
	return s.feishuClient
}

// Shutdown 关闭服务器
func (s *Server) Shutdown() {
	s.logger.Info().Msg("关闭MCP飞书服务器")
}
