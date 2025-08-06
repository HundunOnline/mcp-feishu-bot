# MCP Feishu Bot - Claude Desktop 扩展

这是一个 Claude Desktop 的 MCP (Model Context Protocol) 扩展，让您可以直接在 Claude 对话中发送各种类型的飞书消息。

## ✨ 功能特性

- 🔧 **简单安装** - 一键安装脚本，自动配置
- 💬 **多种消息类型** - 支持文本、富文本、图片、交互式卡片、群聊分享
- 🔒 **三种安全模式** - 无验证、签名验证、关键词验证
- 🌍 **跨平台支持** - macOS、Linux、Windows
- 📱 **实时对话** - 在 Claude 中直接发送飞书消息

## 🚀 快速安装

### 方法 1：一键安装脚本

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/HundunOnline/mcp-feishu-bot/main/claude-extension/install.sh | bash
```

**Windows:**
1. 下载 [install.bat](https://raw.githubusercontent.com/HundunOnline/mcp-feishu-bot/main/claude-extension/install.bat)
2. 右键选择"以管理员身份运行"

### 方法 2：手动安装

#### 1. 下载二进制文件

从 [Releases](https://github.com/HundunOnline/mcp-feishu-bot/releases) 页面下载适合您系统的版本：

- **macOS**: `mcp-feishu-darwin-amd64.tar.gz` 或 `mcp-feishu-darwin-arm64.tar.gz`
- **Linux**: `mcp-feishu-linux-amd64.tar.gz` 或 `mcp-feishu-linux-arm64.tar.gz`  
- **Windows**: `mcp-feishu-windows-amd64.zip`

#### 2. 安装到系统PATH

解压并将 `mcp-feishu` 可执行文件复制到系统PATH中的目录：

**macOS/Linux:**
```bash
sudo mv mcp-feishu /usr/local/bin/
sudo chmod +x /usr/local/bin/mcp-feishu
```

**Windows:**
将 `mcp-feishu.exe` 复制到 `C:\Program Files\MCP-Feishu\` 并添加到系统PATH。

#### 3. 配置 Claude Desktop

编辑 Claude Desktop 配置文件：

**macOS:**
```bash
nano ~/Library/Application\ Support/Claude/claude_desktop_config.json
```

**Linux:**
```bash
nano ~/.config/claude-desktop/claude_desktop_config.json
```

**Windows:**
```
%APPDATA%\Claude\claude_desktop_config.json
```

添加以下配置：

```json
{
  "mcpServers": {
    "mcp-feishu": {
      "command": "mcp-feishu",
      "args": ["-env"],
      "env": {
        "FEISHU_WEBHOOK_URL": "你的飞书Webhook URL",
        "FEISHU_SECURITY_TYPE": "none",
        "FEISHU_SECRET": "",
        "FEISHU_KEYWORDS": ""
      }
    }
  }
}
```

## ⚙️ 配置设置

### 必需配置

| 配置项 | 说明 | 示例 |
|--------|------|------|
| `FEISHU_WEBHOOK_URL` | 飞书机器人Webhook URL | `https://open.feishu.cn/open-apis/bot/v2/hook/xxx` |

### 可选配置

| 配置项 | 说明 | 可选值 | 默认值 |
|--------|------|--------|--------|
| `FEISHU_SECURITY_TYPE` | 安全验证类型 | `none`, `signature`, `keyword` | `none` |
| `FEISHU_SECRET` | 签名验证密钥 | 字符串 | 空 |
| `FEISHU_KEYWORDS` | 关键词列表 | JSON数组字符串 | 空 |

### 安全模式详解

#### 1. 无安全验证 (none)
```json
{
  "FEISHU_SECURITY_TYPE": "none"
}
```

#### 2. 签名验证 (signature)
```json
{
  "FEISHU_SECURITY_TYPE": "signature",
  "FEISHU_SECRET": "your-secret-key"
}
```

#### 3. 关键词验证 (keyword)
```json
{
  "FEISHU_SECURITY_TYPE": "keyword",
  "FEISHU_KEYWORDS": "[\"关键词1\", \"关键词2\", \"keyword\"]"
}
```

## 🎯 使用方法

安装并配置完成后，重启 Claude Desktop，然后您可以在对话中使用以下功能：

### 发送文本消息
```
请发送一条测试消息到飞书："Hello from Claude!"
```

### 发送富文本消息
```
请发送一条包含链接的富文本消息到飞书，内容是访问飞书官网了解更多信息
```

### 发送交互式卡片
```
请发送一个任务完成通知卡片到飞书，包含确认和查看详情按钮
```

### 发送图片消息
```
请发送一张图片到飞书，image_key是：img_v2_041b28e3-xxxx-xxxx-xxxx-xxxxxxxxxxxx
```

## 🛠 故障排除

### 1. 命令未找到错误

**问题**: `mcp-feishu: command not found`

**解决方案**:
- 确保二进制文件在系统PATH中
- 检查文件是否有执行权限
- 在配置中使用完整路径：`"command": "/usr/local/bin/mcp-feishu"`

### 2. 权限错误

**问题**: Permission denied

**解决方案**:
```bash
sudo chmod +x /usr/local/bin/mcp-feishu
```

### 3. 飞书消息发送失败

**问题**: 消息发送到飞书失败

**解决方案**:
- 检查 `FEISHU_WEBHOOK_URL` 是否正确
- 确认飞书机器人配置正确
- 检查网络连接
- 验证安全设置（如果使用签名或关键词模式）

### 4. Claude Desktop 无法识别扩展

**问题**: Claude 中看不到MCP工具

**解决方案**:
- 检查配置文件格式是否正确（JSON语法）
- 重启 Claude Desktop 应用
- 查看 Claude Desktop 日志文件

### 5. 调试模式

启用调试日志：
```json
{
  "mcpServers": {
    "mcp-feishu": {
      "command": "mcp-feishu",
      "args": ["-env", "-debug"],
      "env": {
        "FEISHU_WEBHOOK_URL": "your-webhook-url"
      }
    }
  }
}
```

## 📚 示例对话

以下是一些您可以在 Claude 中尝试的示例：

### 基础文本消息
```
用户: 发送消息"系统维护通知：今晚22:00-24:00进行维护"到飞书
Claude: 我来为您发送这条维护通知到飞书...
[使用 send_text_message 工具]
```

### 富文本消息
```
用户: 发送一条包含链接的消息到飞书，告诉团队访问飞书官网了解新功能
Claude: 我来发送一条包含链接的富文本消息...
[使用 send_post_message 工具]
```

### 交互式卡片
```
用户: 发送一个项目状态报告卡片到飞书，包含项目名称、完成度和操作按钮
Claude: 我来创建一个项目状态报告卡片...
[使用 send_interactive_message 工具]
```

## 🔗 相关链接

- [项目主页](https://github.com/HundunOnline/mcp-feishu-bot)
- [完整文档](https://github.com/HundunOnline/mcp-feishu-bot/blob/main/README.md)
- [问题反馈](https://github.com/HundunOnline/mcp-feishu-bot/issues)
- [飞书开放平台文档](https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot)
- [MCP协议文档](https://modelcontextprotocol.io/)

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](../LICENSE) 文件了解详情。
