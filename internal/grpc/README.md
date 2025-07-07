# og

Go gRPC项目快速开发脚手架

- ORM使用 [ent](https://github.com/ent/ent)
- Redis使用 [go-redis](https://github.com/redis/go-redis)
- 日志使用 [zap](https://github.com/uber-go/zap)
- 配置使用 [viper](https://github.com/spf13/viper)
- 命令行使用 [cli](https://github.com/urfave/cli)
- 工具包使用 [og](https://github.com/noble-gase/ne)
- 使用 [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) 同时支持 grpc 和 http 服务
- http服务支持跨域
- 支持 proto 参数验证
- 支持 swagger.json 生成
- 包含 TraceId、请求日志 等中间件
- 简单好用的 Result Status 统一输出方式

### 前提条件

```shell
# 安装依赖工具
sh install.sh

# 生成ent模块
og ent --help
```
