# og

Go gRPC项目快速开发脚手架

- 数据库使用 [jet](https://github.com/go-jet/jet)
- Redis使用 [go-redis](https://github.com/redis/go-redis)
- 配置使用 [viper](https://github.com/spf13/viper)
- 命令行使用 [cli](https://github.com/urfave/cli)
- 工具包使用 [og](https://github.com/noble-gase/ne)
- 使用 [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) 支持 HTTP 服务
- HTTP 服务支持跨域
- 支持 proto 参数验证
- 支持 swagger.json 生成
- 包含 TraceId、请求日志 等中间件
- 简单好用的 Result Status 统一输出方式

#### 前提条件

```shell
# 安装 protoc
brew install protobuf

# 安装依赖工具
sh install.sh
```

#### Ent支持

```shell
# 安装 ent
go install entgo.io/ent/cmd/ent@latest

# 生成 ent 模块
og ent --help
```
