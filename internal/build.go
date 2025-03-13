package internal

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/noble-gase/og/internal/ent"
	"github.com/noble-gase/og/internal/grpc"
	"github.com/noble-gase/og/internal/http"
)

type Params struct {
	Module  string
	AppPkg  string
	AppName string
	DockerF string
	Proto   bool
}

func InitHttpProject(root, mod string, proto bool, apps ...string) {
	// 创建项目
	initProject(root, mod, proto, http.FS)
	// 创建App(单应用)
	if len(apps) == 0 {
		initApp(root, mod, "", proto, http.FS)
		return
	}
	// 创建App(多应用)
	for _, name := range apps {
		initApp(root, mod, name, proto, http.FS)
	}
}

func InitHttpApp(root, mod, name string, proto bool) {
	initApp(root, mod, name, proto, http.FS)
}

func InitGrpcProject(root, mod string, apps ...string) {
	// 创建项目
	initProject(root, mod, false, grpc.FS)
	// 创建App(单应用)
	if len(apps) == 0 {
		initApp(root, mod, "", false, grpc.FS)
		return
	}
	// 创建App(多应用)
	for _, name := range apps {
		initApp(root, mod, name, false, grpc.FS)
	}
}

func InitGrpcApp(root, mod, name string) {
	initApp(root, mod, name, false, grpc.FS)
}

func InitEnt(root, mod, name string) {
	params := &Params{
		Module:  mod,
		AppName: "ent",
	}
	if len(name) != 0 {
		params.AppName = name
	}
	// ent目录文件
	_ = fs.WalkDir(ent.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) == ".go" {
			return nil
		}
		output := genOutput(root+"/pkg/ent", path, "")
		if len(name) != 0 {
			output = strings.Replace(output, "/ent", "/ent/"+name, 1)
		}
		buildTmpl(ent.FS, path, filepath.Clean(output), params)
		return nil
	})
}

func initProject(root, mod string, proto bool, fsys embed.FS) {
	params := &Params{
		Module: mod,
		Proto:  proto,
	}
	// 项目根目录文件
	files, _ := fs.ReadDir(fsys, ".")
	for _, v := range files {
		if v.IsDir() || filepath.Ext(v.Name()) == ".go" {
			continue
		}
		output := genOutput(root, v.Name(), "")
		buildTmpl(fsys, v.Name(), filepath.Clean(output), params)
	}
	// lib目录文件
	_ = fs.WalkDir(fsys, "pkg/lib", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) == ".go" {
			return nil
		}
		output := genOutput(root, path, "")
		buildTmpl(fsys, path, filepath.Clean(output), params)
		return nil
	})
	// internal目录文件
	_ = fs.WalkDir(fsys, "pkg/internal", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) == ".go" {
			return nil
		}
		// proto模式下，无需code
		if proto && path == "pkg/internal/code/code.tmpl" {
			return nil
		}
		output := genOutput(root, path, "")
		buildTmpl(fsys, path, filepath.Clean(output), params)
		return nil
	})
}

func initApp(root, mod, name string, proto bool, fsys embed.FS) {
	params := &Params{
		Module:  mod,
		AppPkg:  "app",
		AppName: root,
		DockerF: "Dockerfile",
		Proto:   proto,
	}
	if len(name) != 0 {
		params.AppPkg = "app/" + name
		params.AppName = name
		params.DockerF = name + ".dockerfile"
	}
	// app目录文件
	if proto {
		_ = fs.WalkDir(fsys, "pkg/proto", func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() || filepath.Ext(path) == ".go" {
				return nil
			}
			appPath := strings.Replace(path, "proto", "app", 1)
			output := genOutput(root, appPath, name)
			if len(name) != 0 {
				output = strings.Replace(output, "/app", "/app/"+name, 1)
			}
			buildTmpl(fsys, path, filepath.Clean(output), params)
			return nil
		})
		return
	}
	_ = fs.WalkDir(fsys, "pkg/app", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) == ".go" {
			return nil
		}
		output := genOutput(root, path, name)
		if len(name) != 0 {
			output = strings.Replace(output, "/app", "/app/"+name, 1)
		}
		buildTmpl(fsys, path, filepath.Clean(output), params)
		return nil
	})
}

func genOutput(root, path, appName string) string {
	var builder strings.Builder
	// 项目根目录
	builder.WriteString(root)
	builder.WriteString("/")
	// 解析path
	dir, name := filepath.Split(path)
	// dockerfile
	switch name {
	case "Dockerfile":
		if len(appName) != 0 {
			builder.WriteString(appName)
			builder.WriteString(".dockerfile")
		} else {
			builder.WriteString("Dockerfile")
		}
		return filepath.Clean(builder.String())
	case "dockerun.sh":
		if len(appName) != 0 {
			builder.WriteString(appName)
			builder.WriteString("_dockerun.sh")
		} else {
			builder.WriteString("dockerun.sh")
		}
		return filepath.Clean(builder.String())
	}
	// 文件目录
	if len(dir) != 0 {
		builder.WriteString(dir)
	}
	// 文件名称
	switch ext := filepath.Ext(path); ext {
	case ".tmpl":
		builder.WriteString(name[:len(name)-5])
		builder.WriteString(".go")
	case "":
		if strings.Contains(name, "ignore") {
			builder.WriteString(".")
		}
		builder.WriteString(name)
	default:
		builder.WriteString(name)
	}
	// 文件路径
	return builder.String()
}

func buildTmpl(fsys embed.FS, path, output string, params *Params) {
	b, err := fsys.ReadFile(path)
	if err != nil {
		log.Fatalln(FmtErr(err))
	}
	// 模板解析
	t, err := template.New(path).Parse(string(b))
	if err != nil {
		log.Fatalln(FmtErr(err))
	}
	// 文件创建
	f, err := CreateFile(output)
	if err != nil {
		log.Fatalln(FmtErr(err))
	}
	defer f.Close()
	// 模板执行
	if err = t.Execute(f, &params); err != nil {
		log.Fatalln(FmtErr(err))
	}
	fmt.Println(output)
}

func IsDirEmpty(path string) (string, bool) {
	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		log.Fatalln(FmtErr(err))
	}

	// Open the directory
	dir, err := os.Open(absPath)
	if err != nil {
		if os.IsNotExist(err) {
			return absPath, true
		}
		log.Fatalln(FmtErr(err))
	}
	defer dir.Close()

	// Read directory contents
	entries, err := dir.Readdirnames(1) // Read at most one entry
	if err != nil {
		if err == io.EOF {
			return absPath, true // Directory is empty
		}
		log.Fatalln(FmtErr(err))
	}
	return absPath, len(entries) == 0
}

// CreateFile 创建或清空指定的文件
// 文件已存在，则清空；文件或目录不存在，则以0775权限创建
func CreateFile(filename string) (*os.File, error) {
	abspath, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	if err = os.MkdirAll(path.Dir(abspath), 0o775); err != nil {
		return nil, err
	}
	return os.OpenFile(abspath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o775)
}

func FmtErr(err error) error {
	funcName := ""
	// Skip level 1 to get the caller function
	pc, file, line, _ := runtime.Caller(1)
	// Get the function details
	if fn := runtime.FuncForPC(pc); fn != nil {
		name := fn.Name()
		funcName = name[strings.Index(name, ".")+1:]
	}
	return fmt.Errorf("[%s(%s:%d)] %w", funcName, file, line, err)
}

// CmdExamples formats the given examples to the cli.
func CmdExamples(ex ...string) string {
	for i := range ex {
		ex[i] = "  " + ex[i] // indent each row with 2 spaces.
	}
	return strings.Join(ex, "\n")
}
