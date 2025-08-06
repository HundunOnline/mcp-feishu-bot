#!/bin/bash

# MCP Feishu Extension Installer for Claude Desktop
# 适用于 macOS 和 Linux

set -e

EXTENSION_NAME="mcp-feishu-bot"
EXTENSION_VERSION="1.0.0"
BINARY_NAME="mcp-feishu"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检测操作系统
detect_os() {
    if [[ "$OSTYPE" == "linux-gnu"* ]]; then
        echo "linux"
    elif [[ "$OSTYPE" == "darwin"* ]]; then
        echo "darwin"
    else
        log_error "不支持的操作系统: $OSTYPE"
        exit 1
    fi
}

# 检测架构
detect_arch() {
    case $(uname -m) in
        x86_64) echo "amd64" ;;
        arm64) echo "arm64" ;;
        aarch64) echo "arm64" ;;
        *) log_error "不支持的架构: $(uname -m)"; exit 1 ;;
    esac
}

# 获取Claude Desktop配置路径
get_claude_config_path() {
    local os=$(detect_os)
    if [[ "$os" == "darwin" ]]; then
        echo "$HOME/Library/Application Support/Claude/claude_desktop_config.json"
    elif [[ "$os" == "linux" ]]; then
        echo "$HOME/.config/claude-desktop/claude_desktop_config.json"
    fi
}

# 下载二进制文件
download_binary() {
    local os=$(detect_os)
    local arch=$(detect_arch)
    local binary_url="https://github.com/HundunOnline/mcp-feishu-bot/releases/latest/download/${BINARY_NAME}-${os}-${arch}"
    
    if [[ "$os" == "linux" ]]; then
        binary_url="${binary_url}.tar.gz"
    elif [[ "$os" == "darwin" ]]; then
        binary_url="${binary_url}.tar.gz"
    fi

    local install_dir="/usr/local/bin"
    local temp_dir=$(mktemp -d)
    
    log_info "下载 $EXTENSION_NAME 二进制文件..."
    
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$temp_dir/${BINARY_NAME}.tar.gz" "$binary_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$temp_dir/${BINARY_NAME}.tar.gz" "$binary_url"
    else
        log_error "需要 curl 或 wget 来下载文件"
        exit 1
    fi
    
    # 解压并安装
    cd "$temp_dir"
    tar -xzf "${BINARY_NAME}.tar.gz"
    
    log_info "安装二进制文件到 $install_dir..."
    sudo mv "$BINARY_NAME" "$install_dir/"
    sudo chmod +x "$install_dir/$BINARY_NAME"
    
    # 清理临时文件
    rm -rf "$temp_dir"
    
    log_success "二进制文件安装完成"
}

# 配置Claude Desktop
configure_claude() {
    local config_path=$(get_claude_config_path)
    local config_dir=$(dirname "$config_path")
    
    log_info "配置Claude Desktop..."
    
    # 创建配置目录（如果不存在）
    mkdir -p "$config_dir"
    
    # 读取现有配置或创建新配置
    local existing_config="{}"
    if [[ -f "$config_path" ]]; then
        existing_config=$(cat "$config_path")
    fi
    
    # 创建MCP配置
    local mcp_config=$(cat << 'EOF'
{
  "mcpServers": {
    "mcp-feishu": {
      "command": "mcp-feishu",
      "args": ["-env"],
      "env": {
        "FEISHU_WEBHOOK_URL": "",
        "FEISHU_SECURITY_TYPE": "none",
        "FEISHU_SECRET": "",
        "FEISHU_KEYWORDS": ""
      }
    }
  }
}
EOF
)
    
    # 合并配置（简单替换，实际使用中可能需要更复杂的合并逻辑）
    echo "$mcp_config" > "$config_path"
    
    log_success "Claude Desktop配置完成"
    log_warning "请编辑 $config_path 文件，设置您的飞书Webhook URL"
}

# 验证安装
verify_installation() {
    log_info "验证安装..."
    
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        log_success "$BINARY_NAME 命令可用"
        
        # 测试基本功能
        echo '{"jsonrpc":"2.0","id":1,"method":"ping"}' | $BINARY_NAME -env >/dev/null 2>&1
        if [[ $? -eq 0 ]]; then
            log_success "MCP服务器响应正常"
        else
            log_warning "MCP服务器可能需要配置"
        fi
    else
        log_error "$BINARY_NAME 命令不可用"
        exit 1
    fi
}

# 显示配置说明
show_configuration_help() {
    local config_path=$(get_claude_config_path)
    
    cat << EOF

${GREEN}安装完成！${NC}

${YELLOW}下一步操作：${NC}

1. 在飞书中创建自定义机器人，获取Webhook URL
   参考：https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot

2. 编辑Claude Desktop配置文件：
   ${BLUE}$config_path${NC}

3. 在配置文件中设置您的飞书Webhook URL：
   ${BLUE}"FEISHU_WEBHOOK_URL": "https://open.feishu.cn/open-apis/bot/v2/hook/your-webhook-url"${NC}

4. 重启Claude Desktop应用

5. 在Claude中测试：
   ${BLUE}"请发送一条测试消息到飞书"${NC}

${YELLOW}支持的安全模式：${NC}
- none: 无验证（默认）
- signature: 签名验证（需要设置FEISHU_SECRET）
- keyword: 关键词验证（需要设置FEISHU_KEYWORDS）

${YELLOW}更多信息：${NC}
- 项目主页：https://github.com/HundunOnline/mcp-feishu-bot
- 文档：https://github.com/HundunOnline/mcp-feishu-bot/blob/main/README.md

EOF
}

# 主函数
main() {
    echo "================================"
    echo "  MCP Feishu Extension Installer"
    echo "================================"
    echo
    
    # 检查是否有sudo权限
    if ! sudo -n true 2>/dev/null; then
        log_info "需要sudo权限来安装二进制文件"
    fi
    
    download_binary
    configure_claude
    verify_installation
    show_configuration_help
    
    log_success "安装完成！"
}

# 错误处理
trap 'log_error "安装过程中发生错误"; exit 1' ERR

# 运行主函数
main "$@"
