# App - {{.Module}}

1. 配置文件 `config.toml`
2. 执行 `buf generate` 编译proto文件
3. 执行 `buf dep update` 更新proto依赖
4. 执行 `go run main.go` 运行
5. 执行 `go run main.go -h` 查看命令
6. 根目录执行 `og ent <name>` 生成ent模块
7. 查看API文档 `swagger serve xxx.swagger.json`
