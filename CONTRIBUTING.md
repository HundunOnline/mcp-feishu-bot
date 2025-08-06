# 贡献指南

感谢您对 MCP-Feishu 项目的关注！我们欢迎各种形式的贡献。

## 如何贡献

### 报告问题

如果您发现了 bug 或有功能建议，请：

1. 检查现有的 [Issues](https://github.com/your-username/mcp-feishu/issues) 确保问题未被报告
2. 创建一个新的 Issue，详细描述：
   - 问题的详细描述
   - 重现步骤
   - 预期行为和实际行为
   - 系统环境信息（Go版本、操作系统等）

### 提交代码

1. **Fork 项目**
   ```bash
   git clone https://github.com/your-username/mcp-feishu.git
   cd mcp-feishu
   ```

2. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **开发和测试**
   - 确保代码符合 Go 的编码规范
   - 添加必要的测试
   - 运行现有测试确保不会破坏现有功能
   ```bash
   go test ./...
   go vet ./...
   go fmt ./...
   ```

4. **提交更改**
   ```bash
   git add .
   git commit -m "feat: 添加新功能描述"
   ```

5. **推送到GitHub**
   ```bash
   git push origin feature/your-feature-name
   ```

6. **创建 Pull Request**
   - 提供清晰的 PR 标题和描述
   - 关联相关的 Issues
   - 确保 CI 测试通过

## 代码规范

### Go 代码风格

- 遵循 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- 使用 `go fmt` 格式化代码
- 使用 `go vet` 检查代码问题
- 添加适当的注释，特别是导出的函数和类型

### 提交信息格式

使用 [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) 格式：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

类型包括：
- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档更新
- `style`: 代码格式化
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat(security): 添加签名校验功能

- 实现 HMAC-SHA256 签名验证
- 添加时间戳防重放攻击
- 更新相关文档

Closes #123
```

## 开发环境设置

### 环境要求

- Go 1.21 或更高版本
- Git

### 本地开发

1. **克隆项目**
   ```bash
   git clone https://github.com/your-username/mcp-feishu.git
   cd mcp-feishu
   ```

2. **安装依赖**
   ```bash
   go mod tidy
   ```

3. **配置环境**
   ```bash
   cp examples/env.example .env
   # 编辑 .env 文件，填入您的飞书配置
   ```

4. **运行项目**
   ```bash
   go run main.go
   ```

5. **运行测试**
   ```bash
   go test ./...
   ```

## 项目结构

```
mcp-feishu/
├── main.go                     # 主入口文件
├── internal/                   # 内部包
│   ├── config/                 # 配置管理
│   ├── feishu/                 # 飞书客户端
│   ├── mcp/                    # MCP服务器
│   └── types/                  # 类型定义
├── examples/                   # 配置示例
├── LICENSE                     # 许可证
├── README.md                   # 项目说明
├── CONTRIBUTING.md             # 贡献指南
└── go.mod                      # Go模块定义
```

## 测试

### 单元测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./internal/feishu

# 运行测试并显示覆盖率
go test -cover ./...
```

### 集成测试

在提交 PR 之前，请确保：

1. 配置真实的飞书 Webhook URL
2. 运行示例测试确保功能正常
3. 测试不同的安全设置模式

## 发布流程

1. 更新版本号（在 `main.go` 中）
2. 更新 CHANGELOG.md
3. 创建 Git tag
4. 推送到 GitHub
5. 创建 GitHub Release

## 获取帮助

如果您有任何问题，可以：

- 创建 [Issue](https://github.com/your-username/mcp-feishu/issues)
- 查看项目 [Wiki](https://github.com/your-username/mcp-feishu/wiki)
- 联系维护者

感谢您的贡献！
