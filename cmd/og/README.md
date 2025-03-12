# og

自动生成项目，支持 `HTTP` 和 `gRPC`，并同时支持创建「单应用」和「多应用」

> HTTP 项目配合 `protoc-gen-og`，支持使用 `proto` 定义API

## 安装

```shell
go install github.com/noble-gase/og/cmd/og@latest
```

## 创建项目

#### 单应用

```shell
# http
og new demo
og new demo --proto # 使用proto定义API
og new demo --mod=xxx.com/demo # 指定module名称
.
├── go.mod
├── go.sum
├── Dockerfile
└── pkg
    ├── app
    │   ├── api
    │   ├── cmd
    │   ├── config
    │   ├── middleware
    │   ├── router
    │   ├── service
    │   ├── web
    │   ├── config.toml
    │   └── main.go
    └── internal

# grpc
og new demo --grpc
og new demo --mod=xxx.com/demo --grpc # 指定module名称
.
├── go.mod
├── go.sum
├── Dockerfile
└── pkg
    ├── app
    │   ├── api
    │   │   └── greeter.proto
    │   ├── cmd
    │   ├── config
    │   ├── server
    │   ├── service
    │   ├── buf.gen.yaml
    │   ├── buf.lock
    │   ├── buf.yaml
    │   ├── config.toml
    │   └── main.go
    └── internal
```

#### 多应用

```shell
# http
og new demo --apps=foo,bar
og new demo --apps=foo,bar --proto # 使用proto定义API
og new demo --mod=xxx.com/demo --apps=foo,bar
og new demo --mod=xxx.com/demo --apps=foo --apps=bar
.
├── go.mod
├── go.sum
├── foo.dockerfile
├── bar.dockerfile
└── pkg
    ├── app
    │   ├── foo
    │   │   ├── api
    │   │   ├── cmd
    │   │   ├── config
    │   │   ├── middleware
    │   │   ├── router
    │   │   ├── service
    │   │   ├── web
    │   │   ├── config.toml
    │   │   └── main.go
    │   ├── bar
    │   │   ├── ...
    │   │   └── main.go
    └── internal

# grpc
og new demo --apps=foo,bar --grpc
og new demo --mod=xxx.com/demo --apps=foo,bar --grpc
og new demo --mod=xxx.com/demo --apps=foo --apps=bar --grpc
.
├── go.mod
├── go.sum
├── foo.dockerfile
├── bar.dockerfile
└── pkg
    ├── app
    │   ├── foo
    │   │   ├── api
    │   │   │   └── greeter.proto
    │   │   ├── cmd
    │   │   ├── config
    │   │   ├── server
    │   │   ├── service
    │   │   ├── buf.gen.yaml
    │   │   ├── buf.lock
    │   │   ├── buf.yaml
    │   │   ├── config.toml
    │   │   └── main.go
    │   ├── bar
    │   │   ├── ...
    │   │   └── main.go
    └── internal
```

## 创建应用

```shell
# 多应用项目适用，需在项目根目录执行（即：go.mod所在目录）
og app foo # 创建HTTP应用 -- foo
og app foo --proto # 使用proto定义API
og app foo --grpc # 创建gRPC应用
og app foo bar # 创建两个HTTP应用 -- foo 和 bar
og app foo bar --grpc # 创建两个gRPC应用 -- foo 和 bar
.
├── go.mod
├── go.sum
├── foo.dockerfile
├── bar.dockerfile
└── pkg
    ├── app
    │   ├── foo
    │   │   ├── ...
    │   │   └── main.go
    │   ├── bar
    │   │   ├── ...
    │   │   └── main.go
    └── internal
```

## 创建Ent实例

#### 单实例

```shell
og ent
.
├── go.mod
├── go.sum
└── pkg
    ├── app
    ├── ent
    │   ├── ...
    │   └── schema
    └── internal
```

#### 多实例

```shell
# 创建Ent实例 -- foo 和 bar
og ent foo bar
.
├── go.mod
├── go.sum
└── pkg
    ├── app
    ├── ent
    │   ├── foo
    │   │   ├── ...
    │   │   └── schema
    │   ├── bar
    │   │   ├── ...
    │   │   └── schema
    └── internal
```
