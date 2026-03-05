# 鿫-Og

[![golang](https://img.shields.io/badge/Language-Go-green.svg?style=flat)](https://golang.org)
[![pkg.go.dev](https://img.shields.io/badge/dev-reference-007d9c?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/noble-gase/og)
[![MIT](http://img.shields.io/badge/license-MIT-brightgreen.svg)](http://opensource.org/licenses/MIT)

[鿫-Og] Go开发脚手架

## og

支持创建 `HTTP` 和 `gRPC` 以及 `MCP` 项目，并同时支持创建「单应用」和「多应用」

```shell
go install github.com/noble-gase/og/cmd/og@latest
```

进一步了解 👉 [详情](https://github.com/noble-gase/og/blob/main/cmd/og/README.md)

## protoc-gen-og

使用 `proto` 定义API，基于 [chi](https://github.com/go-chi/chi) 自动生成路由和服务注册

```shell
go install github.com/noble-gase/og/cmd/protoc-gen-og@latest
```

进一步了解 👉 [详情](https://github.com/noble-gase/og/blob/main/cmd/protoc-gen-og/README.md)

## gg

受 `protoc-gen-go` 启发，为结构体生成 `Get` 方法【支持泛型!!!】，避免空指针引起的Panic

```shell
go install github.com/noble-gase/og/cmd/gg@latest
```

进一步了解 👉 [详情](https://github.com/noble-gase/og/blob/main/cmd/gg/README.md)

**Enjoy 😊**
