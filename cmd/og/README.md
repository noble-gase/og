# og - 项目脚手架

自动生成项目，支持 `HTTP` 和 `gRPC`，并同时支持创建「单应用」和「多应用」

> - 项目结构采用[标准布局](https://github.com/golang-standards/project-layout)
> - `HTTP` 配合 `protoc-gen-og`，支持使用 `proto` 定义API

## 安装

```shell
go install github.com/noble-gase/og/cmd/og@latest
```

## 创建项目

### HTTP

##### 单应用

```shell
og new . # 当前目录初始化
og new demo # 创建demo项目
og new demo --mod=xxx.com/demo # 指定module-path
.
├── cmd
│   ├── config.toml
│   └── main.go
├── internal
│   └── app
│       ├── cmd
│       ├── config
│       ├── handler
│       ├── router
│       └── service
├── pkg
│   └── ...
├── Dockerfile
├── dockerun.sh
├── go.mod
├── go.sum
└── README.md
```

##### 多应用

```shell
og new . --app=foo --app=bar # 当前目录初始化
og new demo --app=foo --app=bar # 创建demo项目
og new demo --mod=xxx.com/demo --app=foo --app=bar # 指定module-path
.
├── cmd
│   ├── foo
│   │   ├── config.toml
│   │   └── main.go
│   └── bar
│       ├── config.toml
│       └── main.go
├── internal
│   └── app
│       ├── foo
│       │   ├── cmd
│       │   ├── config
│       │   ├── handler
│       │   ├── router
│       │   └── service
│       └── bar
│           ├── cmd
│           ├── config
│           ├── handler
│           ├── router
│           └── service
├── pkg
│   └── ...
├── bar.dockerfile
├── bar_dockerun.sh
├── foo.dockerfile
├── foo_dockerun.sh
├── go.mod
├── go.sum
└── README.md
```

### HTTP -- 使用proto定义API

##### 单应用

```shell
og new . --proto # 当前目录初始化
og new demo --proto # 创建demo项目
og new demo --mod=xxx.com/demo --proto # 指定module-path
.
├── api
│   ├── code.proto
│   └── greeter.proto
├── cmd
│   ├── config.toml
│   └── main.go
├── internal
│   └── app
│       ├── cmd
│       ├── config
│       ├── router
│       └── service
├── pkg
│   └── ...
├── buf.gen.yaml
├── buf.lock
├── buf.yaml
├── Dockerfile
├── dockerun.sh
├── go.mod
├── go.sum
└── README.md
```

##### 多应用

```shell
og new . --app=foo --app=bar --proto # 当前目录初始化
og new demo --app=foo --app=bar --proto # 创建demo项目
og new demo --mod=xxx.com/demo --app=foo --app=bar --proto # 指定module-path
.
├── api
│   ├── bar
│   │   ├── code.proto
│   │   └── greeter.proto
│   └── foo
│       ├── code.proto
│       └── greeter.proto
├── cmd
│   ├── bar
│   │   ├── config.toml
│   │   └── main.go
│   └── foo
│       ├── config.toml
│       └── main.go
├── internal
│   └── app
│       ├── bar
│       │   ├── cmd
│       │   ├── config
│       │   ├── router
│       │   └── service
│       └── foo
│           ├── cmd
│           ├── config
│           ├── router
│           └── service
├── pkg
│   └── ...
├── buf.gen.yaml
├── buf.lock
├── buf.yaml
├── foo.dockerfile
├── foo_dockerun.sh
├── bar.dockerfile
├── bar_dockerun.sh
├── go.mod
├── go.sum
└── README.md
```

### gRPC

##### 单应用

```shell
og new . --grpc # 当前目录初始化
og new demo --grpc # 创建demo项目
og new demo --mod=xxx.com/demo --grpc # 指定module-path
.
├── api
│   └── greeter.proto
├── cmd
│   ├── config.toml
│   └── main.go
├── internal
│   └── app
│       ├── cmd
│       ├── config
│       ├── server
│       └── service
├── pkg
│   └── ...
├── buf.gen.yaml
├── buf.lock
├── buf.yaml
├── Dockerfile
├── dockerun.sh
├── go.mod
├── go.sum
└── README.md
```

##### 多应用

```shell
og new . --app=foo --app=bar --grpc # 当前目录初始化
og new demo --app=foo --app=bar --grpc # 创建demo项目
og new demo --mod=xxx.com/demo --app=foo --app=bar --grpc # 指定module-path
.
├── api
│   ├── bar
│   │   └── greeter.proto
│   └── foo
│       └── greeter.proto
├── cmd
│   ├── bar
│   │   ├── config.toml
│   │   └── main.go
│   └── foo
│       ├── config.toml
│       └── main.go
├── internal
│   └── app
│       ├── bar
│       │   ├── cmd
│       │   ├── config
│       │   ├── server
│       │   └── service
│       └── foo
│           ├── cmd
│           ├── config
│           ├── server
│           └── service
├── pkg
│   └── ...
├── buf.gen.yaml
├── buf.lock
├── buf.yaml
├── foo.dockerfile
├── foo_dockerun.sh
├── bar.dockerfile
├── bar_dockerun.sh
├── go.mod
├── go.sum
└── README.md
```

## 创建应用

> 多应用项目适用，需在项目根目录执行（即：`go.mod` 所在目录）

```shell
og app foo bar # 创建两个HTTP应用 -- foo 和 bar
og app foo bar --proto # 使用proto定义API -- foo 和 bar
og app foo bar --grpc # 创建两个gRPC应用 -- foo 和 bar
.
├── api
│   ├── bar
│   └── foo
├── cmd
│   ├── bar
│   │   ├── config.toml
│   │   └── main.go
│   └── foo
│       ├── config.toml
│       └── main.go
├── internal
│   └── app
│       ├── bar
│       └── foo
├── pkg
├── foo.dockerfile
├── foo_dockerun.sh
├── bar.dockerfile
├── bar_dockerun.sh
├── go.mod
├── go.sum
└── README.md
```

## 创建Ent实例

#### 单实例

```shell
og ent
.
├── api
├── cmd
├── internal
│   ├── app
│   └── ent
│       ├── schema
│       └── generate.go
├── pkg
├── go.mod
├── go.sum
└── README.md
```

#### 多实例

```shell
og ent foo bar # 创建Ent实例 -- foo 和 bar
.
├── api
├── cmd
├── internal
│   ├── app
│   └── ent
│       ├── foo
│       │   ├── schema
│       │   └── generate.go
│       └── bar
│           ├── schema
│           └── generate.go
├── pkg
├── go.mod
├── go.sum
└── README.md
```
