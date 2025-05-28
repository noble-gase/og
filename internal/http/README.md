# og

Go Web项目快速开发脚手架

- 路由使用 [chi](https://github.com/go-chi/chi)
- ORM使用 [ent](https://github.com/ent/ent)
- Redis使用 [go-redis](https://github.com/redis/go-redis)
- 日志使用 [zap](https://github.com/uber-go/zap)
- 配置使用 [viper](https://github.com/spf13/viper)
- 命令行使用 [cli](https://github.com/urfave/cli)
- 工具包使用 [ne](https://github.com/noble-gase/ne)
- 包含 TraceId、请求日志、跨域 中间件
- 简单好用的 API Result 统一输出方式

### 前提条件

```sh
# ent
go install entgo.io/ent/cmd/ent@latest

# generate
og ent
```
