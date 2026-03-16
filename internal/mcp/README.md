# og

Go MCP服务快速开发脚手架

- 基于 [mcp-go](https://github.com/mark3labs/mcp-go) 构建
- 数据库使用 [jet](https://github.com/go-jet/jet)
- Redis使用 [go-redis](https://github.com/redis/go-redis)
- 配置使用 [viper](https://github.com/spf13/viper)
- 命令行使用 [cli](https://github.com/urfave/cli)
- 工具包使用 [ne](https://github.com/noble-gase/ne)
- 包含 TraceId、Tool调用日志 中间件

### 运行

#### 启动服务

```shell
cd cmd

# Jet Generate
go run main.go jetgen

# Serve MCP
go run main.go serve
```

#### Ent支持

```shell
# 安装 ent
go install entgo.io/ent/cmd/ent@latest

# 生成 ent 模块
og ent --help
```

### 使用

> 1. 安装 Node.js
> 2. 安装 [Claude Desktop](https://claude.com/download)

#### 环境验证

```shell
curl -N -v http://localhost:9000/mcp/demo/sse

sudo chown -R 501:20 "~/.npm"

npx mcp-remote@latest http://localhost:8080/mcp/demo/sse --allow-http --debug
```

#### Claude 配置

> 配置文件 `~/Library/Application Support/Claude/claude_desktop_config.json`

```json
{
    "mcpServers": {
        "demo": {
            "command": "npx",
            "args": [
                "mcp-remote@latest",
                "http://localhost:8080/mcp/demo/sse",
                "--transport",
                "sse-only",
                "--allow-http"
            ]
        }
    }
}
```
