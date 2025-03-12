# protoc-gen-og

使用 `proto` 定义API，基于 [chi](https://github.com/go-chi/chi) 自动生成路由和服务注册

## 安装

```shell
go install github.com/noble-gase/og/cmd/protoc-gen-og@latest
```

## 使用

```shell
# buf.gen.yaml
version: v2
managed:
  enabled: true
  disable:
    - module: buf.build/googleapis/googleapis
      file_option: go_package
    - module: buf.build/bufbuild/protovalidate
      file_option: go_package_prefix
plugins:
  - local: protoc-gen-go
    out: api
    opt: paths=source_relative
  - local: protoc-gen-og
    out: api
    opt: paths=source_relative
```
