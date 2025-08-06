# MCP飞书服务器

这是一个基于Go语言开发的MCP (Model Context Protocol) 服务器，用于发送各种类型的飞书消息。支持飞书自定义机器人的三种安全设置：无安全、签名校验、自定义关键词。

## 功能特性

### 支持的消息类型

- **文本消息** - 发送纯文本消息
- **富文本消息** - 支持格式化文本、链接、@用户等
- **群名片消息** - 发送结构化的群名片
- **图片消息** - 发送图片消息
- **交互式消息卡片** - 发送自定义卡片和按钮

### 支持的安全设置

1. **无安全设置** (`none`) - 不进行任何安全验证
2. **签名校验** (`signature`) - 使用HMAC-SHA256进行签名验证
3. **自定义关键词** (`keyword`) - 消息必须包含指定关键词

## 🚀 Claude Desktop 扩展

**推荐方式：直接在 Claude Desktop 中使用！**

### 一键安装
```bash
# macOS / Linux
curl -fsSL https://raw.githubusercontent.com/HundunOnline/mcp-feishu-bot/main/claude-extension/install.sh | bash

# Windows: 下载并运行 install.bat（需要管理员权限）
```

安装后重启 Claude Desktop，然后您就可以在对话中直接使用：

```
"请发送一条测试消息到飞书：Hello from Claude!"
```

详细说明请查看：[Claude Desktop 扩展文档](claude-extension/README.md)

---

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置飞书机器人

在飞书开放平台创建自定义机器人，获取Webhook URL，参考：https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot

### 3. 配置服务器

**推荐方式：使用环境变量配置**

复制环境变量示例文件：
```bash
cp examples/env.example .env
```

编辑`.env`文件，填入你的飞书Webhook URL：

```bash
# 飞书机器人配置
FEISHU_WEBHOOK_URL=https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url-here
FEISHU_SECURITY_TYPE=none

# 服务器配置
SERVER_HOST=localhost
SERVER_PORT=3000
```

或者直接设置环境变量：
```bash
export FEISHU_WEBHOOK_URL="https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url-here"
export FEISHU_SECURITY_TYPE="none"
```

#### 备用方式：配置文件

如果环境变量不完整，系统会自动尝试从配置文件补充配置：

```bash
cp examples/config.example.json config.json
# 编辑config.json文件
```

配置文件格式：
```json
{
  "feishu": {
    "webhook_url": "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url-here",
    "security_type": "none"
  },
  "server": {
    "port": 3000,
    "host": "localhost"
  }
}
```

### 4. 运行服务器

**默认运行（自动从环境变量和配置文件加载）：**
```bash
go run main.go
```

**仅使用环境变量（忽略配置文件）：**
```bash
go run main.go -env
```

**指定配置文件作为备用：**
```bash
go run main.go -config config.json
```

**启用调试模式：**
```bash
go run main.go -debug
```

**使用.env文件：**
```bash
# 加载.env文件到环境变量
export $(cat .env | xargs)
go run main.go
```

## 配置策略

### 🔄 **配置优先级**
```
环境变量 > 配置文件 > 默认值
```

### 📝 **配置方式选择**

**生产环境（推荐）：仅使用环境变量**
```bash
# 设置环境变量
export FEISHU_WEBHOOK_URL="your-webhook-url" 
export FEISHU_SECURITY_TYPE="signature"
export FEISHU_SECRET="your-secret"

# 运行（不需要配置文件）
go run main.go
```

**开发环境：环境变量 + 配置文件备选**
```bash
# 配置文件作为备选，环境变量覆盖配置文件
go run main.go -config config.json
```

**强制仅使用环境变量：**
```bash
# 完全忽略配置文件
go run main.go -env
```

## 安全配置

支持的环境变量列表：

| 环境变量 | 描述 | 示例值 | 必填 |
|---------|------|--------|------|
| `FEISHU_WEBHOOK_URL` | 飞书Webhook URL | `https://open.feishu.cn/open-apis/bot/v2/hook/xxx` | ✅ |
| `FEISHU_SECURITY_TYPE` | 安全类型 | `none`, `signature`, `keyword` | ❌ (默认: none) |
| `FEISHU_SECRET` | 签名密钥 | `your-secret-key` | ❌ (signature模式必填) |
| `FEISHU_KEYWORDS` | 关键词列表 | `["关键词1", "关键词2"]` | ❌ (keyword模式必填) |
| `SERVER_HOST` | 服务器主机 | `localhost` | ❌ (默认: localhost) |
| `SERVER_PORT` | 服务器端口 | `3000` | ❌ (默认: 3000) |

### 1. 无安全设置

```bash
FEISHU_WEBHOOK_URL=your-webhook-url
FEISHU_SECURITY_TYPE=none
```

或复制示例：
```bash
cp examples/env.none.example .env
```

### 2. 签名校验

```bash
FEISHU_WEBHOOK_URL=your-webhook-url
FEISHU_SECURITY_TYPE=signature
FEISHU_SECRET=your-secret-key
```

或复制示例：
```bash
cp examples/env.signature.example .env
```

### 3. 自定义关键词

```bash
FEISHU_WEBHOOK_URL=your-webhook-url
FEISHU_SECURITY_TYPE=keyword
FEISHU_KEYWORDS='["关键词1", "关键词2", "keyword"]'
```

或复制示例：
```bash
cp examples/env.keyword.example .env
```

## MCP工具列表

支持飞书官方的5种消息类型：

| 工具名称 | 消息类型 | 描述 | 参数 |
|---------|----------|------|------|
| `send_text_message` | `text` | 发送纯文本消息 | `text: string` |
| `send_post_message` | `post` | 发送富文本消息（支持可选标题） | `content: array, title?: string` |
| `send_image_message` | `image` | 发送图片消息 | `image_key: string` |
| `send_interactive_message` | `interactive` | 发送交互式消息卡片 | `elements: array, config?: object, header?: object` |
| `send_share_chat_message` | `share_chat` | 发送群聊分享卡片 | `share_chat_id: string` |

## 使用示例

### 发送文本消息

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "send_text_message",
    "arguments": {
      "text": "你好，这是一条测试消息！"
    }
  }
}
```

### 发送富文本消息

**无标题富文本：**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "send_post_message",
    "arguments": {
      "content": [
        [
          {
            "tag": "text",
            "text": "这是富文本消息，支持"
          },
          {
            "tag": "a",
            "text": "链接",
            "href": "https://www.feishu.cn"
          },
          {
            "tag": "text",
            "text": "和其他格式"
          }
        ]
      ]
    }
  }
}
```

**有标题富文本：**
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "send_post_message",
    "arguments": {
      "title": "📊 数据报告",
      "content": [
        [
          {
            "tag": "text",
            "text": "本月销售额增长15%"
          }
        ]
      ]
    }
  }
}
```

### 发送交互式卡片

```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "send_interactive_message",
    "arguments": {
      "header": {
        "title": {
          "tag": "plain_text",
          "content": "通知卡片"
        }
      },
      "elements": [
        {
          "tag": "div",
          "text": {
            "tag": "plain_text",
            "content": "这是一个通知消息"
          }
        },
        {
          "tag": "button",
          "text": {
            "tag": "plain_text",
            "content": "确认"
          },
          "value": "confirm",
          "type": "primary"
        }
      ]
    }
  }
}
```

## 编译和部署

### 编译二进制文件

```bash
# 编译当前平台
go build -o mcp-feishu main.go

# 交叉编译Linux版本
GOOS=linux GOARCH=amd64 go build -o mcp-feishu-linux main.go

# 交叉编译Windows版本
GOOS=windows GOARCH=amd64 go build -o mcp-feishu.exe main.go
```

### Docker部署

创建`Dockerfile`：

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod tidy && go build -o mcp-feishu main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/mcp-feishu .
COPY --from=builder /app/examples/config.example.json ./config.json
CMD ["./mcp-feishu", "-config", "config.json"]
```

构建和运行：

```bash
docker build -t mcp-feishu .
docker run -e FEISHU_WEBHOOK_URL="your-webhook-url" mcp-feishu
```

## 项目结构

```
mcp-feishu/
├── main.go                     # 主入口文件
├── go.mod                      # Go模块定义
├── internal/                   # 内部包
│   ├── config/                 # 配置管理
│   │   └── config.go
│   ├── feishu/                 # 飞书客户端
│   │   ├── client.go          # HTTP客户端
│   │   ├── message.go         # 消息构建器
│   │   └── security.go        # 安全管理
│   ├── mcp/                   # MCP服务器
│   │   ├── server.go          # 服务器实现
│   │   └── tools.go           # 工具处理
│   └── types/                 # 类型定义
│       └── types.go
└── examples/                  # 配置示例
    ├── config.example.json    # 基础配置
    ├── config.signature.json  # 签名校验配置
    └── config.keyword.json    # 关键词配置
```

## 开发和测试

### 运行测试

```bash
go test ./...
```

### 启用调试日志

```bash
go run main.go -config config.json -debug
```

### 使用curl测试

```bash
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"send_text_message","arguments":{"text":"Hello from curl!"}}}' | go run main.go -config config.json
```

## 错误处理

服务器会返回标准的JSON-RPC错误响应：

- `-32601`: 方法未找到
- `-32602`: 参数无效
- `-32603`: 内部错误

飞书API错误会在工具结果中返回详细信息。

## 贡献

欢迎贡献代码！请查看 [CONTRIBUTING.md](CONTRIBUTING.md) 了解如何参与项目。

## 更新日志

查看 [CHANGELOG.md](CHANGELOG.md) 了解版本更新历史。

## 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 参考资源

- [飞书开放平台文档](https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot)
- [MCP协议规范](https://modelcontextprotocol.io/)
- [Go语言官方文档](https://golang.org/doc/)