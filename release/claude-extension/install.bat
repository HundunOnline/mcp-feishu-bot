@echo off
REM MCP Feishu Extension Installer for Claude Desktop (Windows)

setlocal enabledelayedexpansion

set EXTENSION_NAME=mcp-feishu-bot
set EXTENSION_VERSION=1.0.0
set BINARY_NAME=mcp-feishu

echo ================================
echo   MCP Feishu Extension Installer
echo ================================
echo.

REM 检查管理员权限
net session >nul 2>&1
if %errorLevel% == 0 (
    echo [INFO] 检测到管理员权限
) else (
    echo [ERROR] 需要管理员权限运行此脚本
    echo 请右键点击此文件，选择"以管理员身份运行"
    pause
    exit /b 1
)

REM 检测架构
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set ARCH=amd64
) else if "%PROCESSOR_ARCHITECTURE%"=="ARM64" (
    set ARCH=arm64
) else (
    echo [ERROR] 不支持的架构: %PROCESSOR_ARCHITECTURE%
    pause
    exit /b 1
)

REM 设置下载URL和路径
set BINARY_URL=https://github.com/HundunOnline/mcp-feishu-bot/releases/latest/download/%BINARY_NAME%-windows-%ARCH%.zip
set INSTALL_DIR=%ProgramFiles%\MCP-Feishu
set TEMP_DIR=%TEMP%\mcp-feishu-install

echo [INFO] 下载 %EXTENSION_NAME% 二进制文件...

REM 创建临时目录
if exist "%TEMP_DIR%" rmdir /s /q "%TEMP_DIR%"
mkdir "%TEMP_DIR%"

REM 下载文件
powershell -Command "& {[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12; Invoke-WebRequest -Uri '%BINARY_URL%' -OutFile '%TEMP_DIR%\%BINARY_NAME%.zip'}"
if %errorLevel% neq 0 (
    echo [ERROR] 下载失败
    pause
    exit /b 1
)

REM 解压文件
powershell -Command "Expand-Archive -Path '%TEMP_DIR%\%BINARY_NAME%.zip' -DestinationPath '%TEMP_DIR%' -Force"
if %errorLevel% neq 0 (
    echo [ERROR] 解压失败
    pause
    exit /b 1
)

echo [INFO] 安装二进制文件到 %INSTALL_DIR%...

REM 创建安装目录
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM 复制文件
copy "%TEMP_DIR%\%BINARY_NAME%.exe" "%INSTALL_DIR%\" >nul
if %errorLevel% neq 0 (
    echo [ERROR] 文件复制失败
    pause
    exit /b 1
)

REM 添加到PATH
echo [INFO] 添加到系统PATH...
setx /M PATH "%PATH%;%INSTALL_DIR%"

REM 清理临时文件
rmdir /s /q "%TEMP_DIR%"

echo [INFO] 配置Claude Desktop...

REM 获取Claude Desktop配置路径
set CLAUDE_CONFIG_PATH=%APPDATA%\Claude\claude_desktop_config.json
set CLAUDE_CONFIG_DIR=%APPDATA%\Claude

REM 创建配置目录
if not exist "%CLAUDE_CONFIG_DIR%" mkdir "%CLAUDE_CONFIG_DIR%"

REM 创建配置文件
echo { > "%CLAUDE_CONFIG_PATH%"
echo   "mcpServers": { >> "%CLAUDE_CONFIG_PATH%"
echo     "mcp-feishu": { >> "%CLAUDE_CONFIG_PATH%"
echo       "command": "mcp-feishu.exe", >> "%CLAUDE_CONFIG_PATH%"
echo       "args": ["-env"], >> "%CLAUDE_CONFIG_PATH%"
echo       "env": { >> "%CLAUDE_CONFIG_PATH%"
echo         "FEISHU_WEBHOOK_URL": "", >> "%CLAUDE_CONFIG_PATH%"
echo         "FEISHU_SECURITY_TYPE": "none", >> "%CLAUDE_CONFIG_PATH%"
echo         "FEISHU_SECRET": "", >> "%CLAUDE_CONFIG_PATH%"
echo         "FEISHU_KEYWORDS": "" >> "%CLAUDE_CONFIG_PATH%"
echo       } >> "%CLAUDE_CONFIG_PATH%"
echo     } >> "%CLAUDE_CONFIG_PATH%"
echo   } >> "%CLAUDE_CONFIG_PATH%"
echo } >> "%CLAUDE_CONFIG_PATH%"

echo [SUCCESS] 安装完成！
echo.
echo 下一步操作：
echo.
echo 1. 在飞书中创建自定义机器人，获取Webhook URL
echo    参考：https://open.feishu.cn/document/client-docs/bot-v3/add-custom-bot
echo.
echo 2. 编辑Claude Desktop配置文件：
echo    %CLAUDE_CONFIG_PATH%
echo.
echo 3. 在配置文件中设置您的飞书Webhook URL
echo.
echo 4. 重启Claude Desktop应用
echo.
echo 5. 在Claude中测试："请发送一条测试消息到飞书"
echo.
echo 更多信息：https://github.com/HundunOnline/mcp-feishu-bot
echo.

pause
