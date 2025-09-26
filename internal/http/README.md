# og

Go Web项目快速开发脚手架

- 路由使用 [chi](https://github.com/go-chi/chi)
- 数据库使用 [jet](https://github.com/go-jet/jet)
- Redis使用 [go-redis](https://github.com/redis/go-redis)
- 配置使用 [viper](https://github.com/spf13/viper)
- 命令行使用 [cli](https://github.com/urfave/cli)
- 工具包使用 [ne](https://github.com/noble-gase/ne)
- 包含 TraceId、请求日志、跨域 中间件
- 简单好用的 API Result 统一输出方式

#### 运行

```shell
cd cmd

# Jet Generate
go run main.go jetgen

# Serve HTTP
go run main.go serve
```

#### Ent支持

```shell
# 安装 ent
go install entgo.io/ent/cmd/ent@latest

# 生成 ent 模块
og ent --help
```
