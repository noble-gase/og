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
	"github.com/noble-gase/og/internal/proto"
)

type Params struct {
	Module  string
	ApiPkg  string
	ApiName string
	AppPkg  string
	AppName string
	DockerF string
}

func InitHttpProject(root, mod string, pb bool, apps ...string) {
	fsys := http.FS
	if pb {
		fsys = proto.FS
	}

	params := &Params{
		Module:  mod,
		ApiPkg:  "api",
		ApiName: "api",
		AppPkg:  "app",
		AppName: root,
		DockerF: "Dockerfile",
	}

	genRoot(root, params, fsys)
	genPkg(root, params, fsys)

	if len(apps) == 0 {
		if pb {
			genApi(root, params, fsys)
		}
		genApp(root, params, fsys)
		genCmd(root, params, fsys)
		return
	}

	for _, name := range apps {
		params.ApiPkg = "api/" + name
		params.ApiName = name
		params.AppPkg = "app/" + name
		params.AppName = name
		params.DockerF = "Dockerfile." + name

		if pb {
			genApi(root, params, fsys, name)
		}
		genApp(root, params, fsys, name)
		genCmd(root, params, fsys, name)
	}
}

func InitHttpApp(root, mod, name string, pb bool) {
	fsys := http.FS
	if pb {
		fsys = proto.FS
	}

	params := &Params{
		Module:  mod,
		ApiPkg:  "api/" + name,
		ApiName: name,
		AppPkg:  "app/" + name,
		AppName: name,
		DockerF: "Dockerfile." + name,
	}

	if pb {
		genApi(root, params, fsys, name)
	}
	genApp(root, params, fsys, name)
	genCmd(root, params, fsys, name)
}

func InitGrpcProject(root, mod string, apps ...string) {
	params := &Params{
		Module:  mod,
		ApiPkg:  "api",
		ApiName: "api",
		AppPkg:  "app",
		AppName: root,
		DockerF: "Dockerfile",
	}

	genRoot(root, params, grpc.FS)
	genPkg(root, params, grpc.FS)

	if len(apps) == 0 {
		genApi(root, params, grpc.FS)
		genApp(root, params, grpc.FS)
		genCmd(root, params, grpc.FS)
		return
	}

	for _, name := range apps {
		params.ApiPkg = "api/" + name
		params.ApiName = name
		params.AppPkg = "app/" + name
		params.AppName = name
		params.DockerF = "Dockerfile." + name

		genApi(root, params, grpc.FS, name)
		genApp(root, params, grpc.FS, name)
		genCmd(root, params, grpc.FS, name)
	}
}

func InitGrpcApp(root, mod, name string) {
	params := &Params{
		Module:  mod,
		ApiPkg:  "api/" + name,
		ApiName: name,
		AppPkg:  "app/" + name,
		AppName: name,
		DockerF: "Dockerfile." + name,
	}

	genApi(root, params, grpc.FS, name)
	genApp(root, params, grpc.FS, name)
	genCmd(root, params, grpc.FS, name)
}

func InitEnt(root, mod string, name ...string) {
	params := &Params{
		Module:  mod,
		AppName: "ent",
	}
	if len(name) != 0 {
		params.AppName = name[0]
	}

	_ = fs.WalkDir(ent.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) == ".go" {
			return nil
		}
		output := genOutput(root+"/internal/ent", path, name...)
		if len(name) != 0 {
			output = strings.Replace(output, "/ent", "/ent/"+name[0], 1)
		}
		buildTmpl(ent.FS, path, filepath.Clean(output), params)
		return nil
	})
}

func genRoot(root string, params *Params, fsys embed.FS) {
	files, _ := fs.ReadDir(fsys, ".")
	for _, v := range files {
		if v.IsDir() || filepath.Ext(v.Name()) == ".go" {
			continue
		}
		output := genOutput(root, v.Name(), "")
		buildTmpl(fsys, v.Name(), filepath.Clean(output), params)
	}
}

func genPkg(root string, params *Params, fsys embed.FS) {
	_ = fs.WalkDir(fsys, "pkg", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) == ".go" {
			return nil
		}
		output := genOutput(root, path, "")
		buildTmpl(fsys, path, filepath.Clean(output), params)
		return nil
	})
}

func genApi(root string, params *Params, fsys embed.FS, appname ...string) {
	_ = fs.WalkDir(fsys, "api", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) == ".go" {
			return nil
		}
		output := genOutput(root, path, appname...)
		if len(appname) != 0 {
			output = strings.Replace(output, "/api", "/api/"+appname[0], 1)
		}
		buildTmpl(fsys, path, filepath.Clean(output), params)
		return nil
	})
}

func genApp(root string, params *Params, fsys embed.FS, appname ...string) {
	_ = fs.WalkDir(fsys, "app", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) == ".go" {
			return nil
		}
		output := genOutput(root+"/internal", path, appname...)
		if len(appname) != 0 {
			output = strings.Replace(output, "/app", "/app/"+appname[0], 1)
		}
		buildTmpl(fsys, path, filepath.Clean(output), params)
		return nil
	})
}

func genCmd(root string, params *Params, fsys embed.FS, appname ...string) {
	_ = fs.WalkDir(fsys, "cmd", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || filepath.Ext(path) == ".go" {
			return nil
		}
		output := genOutput(root, path, appname...)
		if len(appname) != 0 {
			output = strings.Replace(output, "/cmd", "/cmd/"+appname[0], 1)
		}
		buildTmpl(fsys, path, filepath.Clean(output), params)
		return nil
	})
}

func genOutput(root, path string, appname ...string) string {
	var builder strings.Builder
	// 根目录
	builder.WriteString(root)
	builder.WriteString("/")
	// 解析path
	dir, name := filepath.Split(path)
	// dockerfile
	switch name {
	case "dockerfile.tmpl":
		if len(appname) != 0 {
			builder.WriteString("Dockerfile.")
			builder.WriteString(appname[0])
		} else {
			builder.WriteString("Dockerfile")
		}
		return filepath.Clean(builder.String())
	case "dockerun.sh":
		if len(appname) != 0 {
			builder.WriteString(appname[0])
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

func GetCurDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return filepath.Base(wd), nil
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
	return fmt.Errorf("🐛 [%s(%s:%d)] %w", funcName, file, line, err)
}

// CmdExamples formats the given examples to the cli.
func CmdExamples(ex ...string) string {
	for i := range ex {
		ex[i] = "  " + ex[i] // indent each row with 2 spaces.
	}
	return strings.Join(ex, "\n")
}
