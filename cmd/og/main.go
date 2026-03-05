package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/noble-gase/og/internal"
	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

func main() {
	cmd := &cobra.Command{
		Use:     "og",
		Short:   "project scaffold",
		Long:    "project scaffold, quickly create a Go project",
		Version: "v0.10.0",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if cmd.Use == "new" && len(args) != 0 {
				if err := os.MkdirAll(args[0], 0o775); err != nil {
					log.Fatalln("🐛 Mkdir failed:", internal.FmtErr(err))
				}
			}
		},
		Run: func(_ *cobra.Command, _ []string) {
			fmt.Println("🐹 Welcome to use noble-gase[Go] scaffold")
		},
	}
	// 注册命令
	cmd.AddCommand(new(), app(), ent())
	// 执行
	if err := cmd.Execute(); err != nil {
		log.Fatalln("🐛 Cmd execute failed:", internal.FmtErr(err))
	}
}

func new() *cobra.Command {
	var grpc bool
	var proto bool
	var mcp bool
	var mod string
	var apps []string
	cmd := &cobra.Command{
		Use:   "new",
		Short: "create a project",
		Example: internal.CmdExamples(
			"👉 -- HTTP --",
			"og new .",
			"og new demo",
			"og new demo --mod xxx.com/demo",
			"og new demo --app foo --app bar",
			"og new demo --mod xxx.com/demo --app foo --app bar",
			"",
			"👉 -- HTTP(proto) --",
			"og new . --proto",
			"og new demo --proto",
			"og new demo --mod xxx.com/demo --proto",
			"og new demo --app foo --app bar --proto",
			"og new demo --mod xxx.com/demo --app foo --app bar --proto",
			"",
			"👉 -- gRPC --",
			"og new . --grpc",
			"og new demo --grpc",
			"og new demo --mod xxx.com/demo --grpc",
			"og new demo --app foo --app bar --grpc",
			"og new demo --mod xxx.com/demo --app foo --app bar --grpc",
			"",
			"👉 -- MCP --",
			"og new . --mcp",
			"og new demo --mcp",
			"og new demo --mod xxx.com/demo --mcp",
			"og new demo --app foo --app bar --mcp",
			"og new demo --mod xxx.com/demo --app foo --app bar --mcp",
		),
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("must specify a project name")
			}
			return nil
		},
		Run: func(_ *cobra.Command, args []string) {
			workDir := args[0]
			if workDir == "." {
				// 判断是否存在go.mod
				_, err := os.Stat("go.mod")
				if err == nil || !os.IsNotExist(err) {
					log.Fatalln("👿 The current directory already exists go.mod, please confirm!")
				}
				if len(mod) == 0 {
					mod, err = internal.GetCurDir()
					if err != nil {
						log.Fatalln("🐛 Failed to get the current directory:", internal.FmtErr(err))
					}
				}
			} else {
				// 判断目录是否为空
				if path, ok := internal.IsDirEmpty(workDir); !ok {
					log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
				}
				if len(mod) == 0 {
					mod = workDir
				}
			}
			// 创建项目文件
			fmt.Println("🐹 Create project files")
			if mcp {
				internal.InitMcpProject(workDir, mod, apps...)
			} else {
				if grpc {
					internal.InitGrpcProject(workDir, mod, apps...)
				} else {
					internal.InitHttpProject(workDir, mod, proto, apps...)
				}
			}
			// go mod init
			fmt.Println("⌛️ go mod init")
			modInit := exec.Command("go", "mod", "init", mod)
			modInit.Dir = workDir
			if err := modInit.Run(); err != nil {
				log.Fatalln("🐛 go mod init failed:", internal.FmtErr(err))
			}
			// go mod tidy
			fmt.Println("⌛️ go mod tidy")
			modTidy := exec.Command("go", "mod", "tidy")
			modTidy.Dir = workDir
			modTidy.Stderr = os.Stderr
			if err := modTidy.Run(); err != nil {
				log.Fatalln("🐛 go mod tidy failed:", internal.FmtErr(err))
			}
			fmt.Println("🐹 Project creation completed! please read README")
		},
	}
	// 注册参数
	cmd.Flags().BoolVar(&grpc, "grpc", false, "create a gRPC project")
	cmd.Flags().BoolVar(&proto, "proto", false, "use proto to define the API")
	cmd.Flags().BoolVar(&mcp, "mcp", false, "create a MCP project")
	cmd.Flags().StringVar(&mod, "mod", "", "set the module name (default is the project name)")
	cmd.Flags().StringSliceVar(&apps, "app", nil, "create a multi-application project")
	return cmd
}

func app() *cobra.Command {
	var grpc bool
	var proto bool
	var mcp bool
	cmd := &cobra.Command{
		Use:   "app",
		Short: "create an application",
		Example: internal.CmdExamples(
			"👉 -- HTTP --",
			"og app foo",
			"og app foo bar",
			"",
			"👉 -- HTTP(proto) --",
			"og app foo --proto",
			"og app foo bar --proto",
			"",
			"👉 -- gRPC --",
			"og app foo --grpc",
			"og app foo bar --grpc",
			"",
			"👉 -- MCP --",
			"og app foo --mcp",
			"og app foo bar --mcp",
		),
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errors.New("must specify an app name")
			}
			return nil
		},
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println("⌛️ Parse go.mod")
			// 读取 go.mod 文件
			data, err := os.ReadFile("go.mod")
			if err != nil {
				log.Fatalln("🐛 Failed to read go.mod file:", internal.FmtErr(err))
			}
			// 解析 go.mod 文件
			f, err := modfile.Parse("go.mod", data, nil)
			if err != nil {
				log.Fatalln("🐛 Failed to parse go.mod file:", internal.FmtErr(err))
			}
			// 创建应用文件
			fmt.Println("🐹 Create application files")
			if mcp {
				for _, name := range args {
					if path, ok := internal.IsDirEmpty("internal/app/" + name); !ok {
						log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
					}
					internal.InitMcpApp(".", f.Module.Mod.Path, name)
				}
			} else {
				if grpc {
					for _, name := range args {
						if path, ok := internal.IsDirEmpty("internal/app/" + name); !ok {
							log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
						}
						internal.InitGrpcApp(".", f.Module.Mod.Path, name)
					}
				} else {
					for _, name := range args {
						if path, ok := internal.IsDirEmpty("internal/app/" + name); !ok {
							log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
						}
						internal.InitHttpApp(".", f.Module.Mod.Path, name, proto)
					}
				}
			}
			// go mod tidy
			fmt.Println("⌛️ go mod tidy")
			modTidy := exec.Command("go", "mod", "tidy")
			modTidy.Stderr = os.Stderr
			if err := modTidy.Run(); err != nil {
				log.Fatalln("🐛 go mod tidy failed:", internal.FmtErr(err))
			}
			fmt.Println("🐹 Application creation completed! please read README")
		},
	}
	// 注册参数
	cmd.Flags().BoolVar(&grpc, "grpc", false, "create a gRPC application")
	cmd.Flags().BoolVar(&proto, "proto", false, "use proto to define the API")
	cmd.Flags().BoolVar(&mcp, "mcp", false, "create a MCP application")
	return cmd
}

func ent() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ent",
		Short: "create an ent module",
		Example: internal.CmdExamples(
			"👉 -- default instance --",
			"og ent",
			"",
			"👉 -- specify name --",
			"og ent foo",
			"og ent foo bar",
		),
		Run: func(_ *cobra.Command, args []string) {
			fmt.Println("⌛️ Parse go.mod")
			// 读取 go.mod 文件
			data, err := os.ReadFile("go.mod")
			if err != nil {
				log.Fatalln("🐛 Failed to read go.mod file:", internal.FmtErr(err))
			}
			// 解析 go.mod 文件
			f, err := modfile.Parse("go.mod", data, nil)
			if err != nil {
				log.Fatalln("🐛 Failed to parse go.mod file:", internal.FmtErr(err))
			}
			// 创建Ent文件
			fmt.Println("🐹 Create ent file")
			if len(args) != 0 {
				for _, name := range args {
					if path, ok := internal.IsDirEmpty("internal/ent/" + name); !ok {
						log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
					}
					internal.InitEnt(".", f.Module.Mod.Path, name)
				}
			} else {
				if path, ok := internal.IsDirEmpty("internal/ent"); !ok {
					log.Fatalf("👿 The directory(%s) is not empty, please confirm!", path)
				}
				internal.InitEnt(".", f.Module.Mod.Path)
			}
			// go mod tidy
			fmt.Println("⌛️ go mod tidy")
			modTidy := exec.Command("go", "mod", "tidy")
			modTidy.Stderr = os.Stderr
			if err := modTidy.Run(); err != nil {
				log.Fatalln("🐛 go mod tidy failed:", internal.FmtErr(err))
			}
			// ent generate
			fmt.Println("⌛️ Ent generate")
			if len(args) != 0 {
				for _, name := range args {
					entGen := exec.Command("go", "generate", "./internal/ent/"+name)
					if err := entGen.Run(); err != nil {
						log.Fatalln("🐛 Ent generate failed:", internal.FmtErr(err))
					}
				}
			} else {
				entGen := exec.Command("go", "generate", "./internal/ent")
				if err := entGen.Run(); err != nil {
					log.Fatalln("🐛 Ent generate failed:", internal.FmtErr(err))
				}
			}
			// go mod tidy
			fmt.Println("⌛️ go mod tidy")
			modClean := exec.Command("go", "mod", "tidy")
			modClean.Stderr = os.Stderr
			if err := modClean.Run(); err != nil {
				log.Fatalln("🐛 go mod tidy failed:", internal.FmtErr(err))
			}
			fmt.Println("🐹 Ent module creation completed! please read README")
		},
	}
	return cmd
}
