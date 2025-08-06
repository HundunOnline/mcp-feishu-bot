// Package main implements a Model Context Protocol (MCP) server for Feishu messaging.
//
// This server provides tools for sending various types of Feishu messages including:
// - Text messages
// - Rich text messages (post format)
// - Image messages
// - Interactive message cards
// - Share chat messages
//
// It supports three security modes:
// - None: No security validation
// - Signature: HMAC-SHA256 signature verification
// - Keyword: Message content must contain specified keywords
//
// The server can be configured via environment variables or JSON configuration files,
// with environment variables taking precedence.
package main

import (
	"flag"
	"fmt"
	"mcp-feishu/internal/config"
	"mcp-feishu/internal/feishu"
	"mcp-feishu/internal/mcp"
	"os"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// 解析命令行参数
	var (
		configPath = flag.String("config", "", "配置文件路径（作为环境变量的备用选项）")
		useEnv     = flag.Bool("env", true, "强制仅使用环境变量（忽略配置文件）")
		debug      = flag.Bool("debug", false, "启用调试日志")
		version    = flag.Bool("version", false, "显示版本信息")
	)
	flag.Parse()

	// 显示版本信息
	if *version {
		fmt.Println("MCP飞书服务器 v1.0.0")
		fmt.Println("支持发送各种飞书消息类型和三种安全设置")
		return
	}

	// 配置日志级别
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// 配置日志输出格式
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: "15:04:05",
	})

	log.Info().Msg("启动MCP飞书服务器")

	// 加载配置 - 优先从环境变量读取
	var cfg *config.Config

	// 首先尝试从环境变量加载配置
	cfg = config.LoadFromEnv()

	// 如果环境变量配置不完整且提供了配置文件，则从配置文件补充
	if cfg.Feishu.WebhookURL == "" && (*configPath != "" || !*useEnv) {
		log.Info().Str("config_path", *configPath).Msg("环境变量配置不完整，尝试从配置文件加载")
		fileCfg, err := config.LoadConfig(*configPath)
		if err != nil {
			log.Fatal().Err(err).Msg("加载配置文件失败")
		}

		// 用配置文件补充环境变量中缺失的配置
		cfg = config.MergeConfig(cfg, fileCfg)
		log.Info().Msg("已合并环境变量和配置文件设置")
	} else {
		log.Info().Msg("从环境变量加载配置")
	}

	log.Info().
		Str("webhook_url", cfg.Feishu.WebhookURL).
		Str("security_type", cfg.Feishu.SecurityType).
		Msg("配置加载成功")

	// 创建飞书客户端
	feishuClient := feishu.NewClient(cfg.Feishu)
	log.Info().Msg("飞书客户端创建成功")

	// 创建MCP服务器
	mcpServer := mcp.NewServer(feishuClient)
	log.Info().Msg("MCP服务器创建成功")

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		if err := mcpServer.Run(); err != nil {
			log.Fatal().Err(err).Msg("MCP服务器运行失败")
		}
	}()

	log.Info().Msg("MCP飞书服务器启动成功，等待请求...")

	// 等待退出信号
	sig := <-sigChan
	log.Info().Str("signal", sig.String()).Msg("收到退出信号")

	// 优雅关闭
	mcpServer.Shutdown()
	log.Info().Msg("MCP飞书服务器已关闭")
}
