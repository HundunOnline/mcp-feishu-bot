# 飞书消息测试示例

这里提供了各种类型的飞书消息测试示例，可以用于测试MCP服务器的功能。

## 1. 文本消息

### 简单文本
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "tools/call",
  "params": {
    "name": "send_text_message",
    "arguments": {
      "text": "这是一条简单的文本消息"
    }
  }
}
```

### 包含关键词的文本（用于关键词安全模式）
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "method": "tools/call",
  "params": {
    "name": "send_text_message",
    "arguments": {
      "text": "这条消息包含关键词1，符合安全要求"
    }
  }
}
```

## 2. 富文本消息

### 基础富文本
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "method": "tools/call",
  "params": {
    "name": "send_post_message",
    "arguments": {
      "content": {
        "post": {
          "zh_cn": {
            "content": [
              [
                {
                  "tag": "text",
                  "text": "这是富文本消息，支持"
                },
                {
                  "tag": "text",
                  "text": "加粗",
                  "style": ["bold"]
                },
                {
                  "tag": "text",
                  "text": "、"
                },
                {
                  "tag": "text",
                  "text": "斜体",
                  "style": ["italic"]
                },
                {
                  "tag": "text",
                  "text": "等格式"
                }
              ]
            ]
          }
        }
      }
    }
  }
}
```

### 包含链接的富文本
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "method": "tools/call",
  "params": {
    "name": "send_post_message",
    "arguments": {
      "content": {
        "post": {
          "zh_cn": {
            "content": [
              [
                {
                  "tag": "text",
                  "text": "访问"
                },
                {
                  "tag": "a",
                  "text": "飞书官网",
                  "href": "https://www.feishu.cn"
                },
                {
                  "tag": "text",
                  "text": "了解更多信息"
                }
              ]
            ]
          }
        }
      }
    }
  }
}
```

### 包含@用户的富文本
```json
{
  "jsonrpc": "2.0",
  "id": 5,
  "method": "tools/call",
  "params": {
    "name": "send_post_message",
    "arguments": {
      "content": {
        "post": {
          "zh_cn": {
            "content": [
              [
                {
                  "tag": "at",
                  "user_id": "all",
                  "user_name": "所有人"
                },
                {
                  "tag": "text",
                  "text": " 请注意这条重要消息！"
                }
              ]
            ]
          }
        }
      }
    }
  }
}
```

## 3. 交互式消息卡片

### 简单信息卡片
```json
{
  "jsonrpc": "2.0",
  "id": 6,
  "method": "tools/call",
  "params": {
    "name": "send_interactive_message",
    "arguments": {
      "header": {
        "title": {
          "tag": "plain_text",
          "content": "系统通知"
        },
        "template": "blue"
      },
      "elements": [
        {
          "tag": "div",
          "text": {
            "tag": "plain_text",
            "content": "这是一条系统通知消息"
          }
        }
      ]
    }
  }
}
```

### 包含按钮的卡片
```json
{
  "jsonrpc": "2.0",
  "id": 7,
  "method": "tools/call",
  "params": {
    "name": "send_interactive_message",
    "arguments": {
      "header": {
        "title": {
          "tag": "plain_text",
          "content": "操作确认"
        },
        "template": "orange"
      },
      "elements": [
        {
          "tag": "div",
          "text": {
            "tag": "plain_text",
            "content": "请确认是否执行此操作？"
          }
        },
        {
          "tag": "action",
          "actions": [
            {
              "tag": "button",
              "text": {
                "tag": "plain_text",
                "content": "确认"
              },
              "value": "confirm",
              "type": "primary"
            },
            {
              "tag": "button",
              "text": {
                "tag": "plain_text",
                "content": "取消"
              },
              "value": "cancel",
              "type": "default"
            }
          ]
        }
      ]
    }
  }
}
```

### 复杂信息卡片
```json
{
  "jsonrpc": "2.0",
  "id": 8,
  "method": "tools/call",
  "params": {
    "name": "send_interactive_message",
    "arguments": {
      "header": {
        "title": {
          "tag": "plain_text",
          "content": "任务状态报告"
        },
        "subtitle": {
          "tag": "plain_text",
          "content": "定时任务执行情况"
        },
        "template": "green"
      },
      "elements": [
        {
          "tag": "div",
          "fields": [
            {
              "is_short": true,
              "text": {
                "tag": "lark_md",
                "content": "**任务名称**\\n数据同步任务"
              }
            },
            {
              "is_short": true,
              "text": {
                "tag": "lark_md",
                "content": "**执行状态**\\n✅ 成功"
              }
            },
            {
              "is_short": true,
              "text": {
                "tag": "lark_md",
                "content": "**执行时间**\\n2024-01-15 10:30:00"
              }
            },
            {
              "is_short": true,
              "text": {
                "tag": "lark_md",
                "content": "**处理记录**\\n1,234 条"
              }
            }
          ]
        },
        {
          "tag": "hr"
        },
        {
          "tag": "div",
          "text": {
            "tag": "lark_md",
            "content": "任务执行完成，所有数据已成功同步到目标系统。"
          }
        },
        {
          "tag": "action",
          "actions": [
            {
              "tag": "button",
              "text": {
                "tag": "plain_text",
                "content": "查看详情"
              },
              "value": "view_details",
              "type": "primary"
            },
            {
              "tag": "button",
              "text": {
                "tag": "plain_text",
                "content": "下载日志"
              },
              "value": "download_log",
              "type": "default"
            }
          ]
        }
      ]
    }
  }
}
```

## 4. 图片消息

```json
{
  "jsonrpc": "2.0",
  "id": 9,
  "method": "tools/call",
  "params": {
    "name": "send_image_message",
    "arguments": {
      "image_key": "img_v2_041b28e3-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
    }
  }
}
```

注意：图片消息需要先上传图片到飞书获取image_key。

## 5. 群名片消息

```json
{
  "jsonrpc": "2.0",
  "id": 10,
  "method": "tools/call",
  "params": {
    "name": "send_share_chat_message",
    "arguments": {
      "share_chat_id": "oc_xxxxxxxxxxxxxxxxxxxxxxxx"
    }
  }
}
```

## 测试命令

### 使用echo和管道测试

```bash
# 测试文本消息
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"send_text_message","arguments":{"text":"Hello World!"}}}' | go run main.go -config config.json

# 测试富文本消息
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"send_rich_text_message","arguments":{"content":{"post":{"zh_cn":{"content":[[{"tag":"text","text":"这是"},{"tag":"a","text":"链接","href":"https://www.feishu.cn"}]]}}}}}}' | go run main.go -config config.json
```

### 使用curl测试（如果服务器运行在HTTP模式）

```bash
curl -X POST http://localhost:3000 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"send_text_message","arguments":{"text":"Hello from curl!"}}}'
```

### 批量测试脚本

创建`test.sh`脚本：

```bash
#!/bin/bash

# 测试文本消息
echo "测试文本消息..."
echo '{"jsonrpc":"2.0","id":1,"method":"tools/call","params":{"name":"send_text_message","arguments":{"text":"测试文本消息"}}}' | go run main.go -config config.json

sleep 2

# 测试富文本消息
echo "测试富文本消息..."
echo '{"jsonrpc":"2.0","id":2,"method":"tools/call","params":{"name":"send_rich_text_message","arguments":{"content":{"post":{"zh_cn":{"content":[[{"tag":"text","text":"这是富文本消息，包含"},{"tag":"a","text":"链接","href":"https://www.feishu.cn"}]]}}}}}}' | go run main.go -config config.json

sleep 2

# 测试卡片消息
echo "测试卡片消息..."
echo '{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"send_interactive_message","arguments":{"header":{"title":{"tag":"plain_text","content":"测试卡片"}},"elements":[{"tag":"div","text":{"tag":"plain_text","content":"这是测试卡片消息"}}]}}}' | go run main.go -config config.json

echo "测试完成！"
```

运行测试：
```bash
chmod +x test.sh
./test.sh
```