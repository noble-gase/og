package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/genproto/googleapis/api/annotations"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	version = "v0.0.2"
	suffix  = "_http.pb.go"
)

func main() {
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Printf("protoc-gen-og %s\n", version)
		return
	}

	var flags flag.FlagSet

	protogen.Options{
		ParamFunc: flags.Set,
	}.Run(func(p *protogen.Plugin) error {
		p.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range p.Files {
			if !f.Generate {
				continue
			}
			genServiceFile(p, f)
			genCodeFile(p, f)
		}
		return nil
	})
}

const (
	ctxPkg     = protogen.GoImportPath("context")
	errorPkg   = protogen.GoImportPath("errors")
	httpPkg    = protogen.GoImportPath("net/http")
	chiPkg     = protogen.GoImportPath("github.com/go-chi/chi/v5")
	contribPkg = protogen.GoImportPath("github.com/noble-gase/ne")
	resultPkg  = protogen.GoImportPath("github.com/noble-gase/ne/result")
	restyPkg   = protogen.GoImportPath("github.com/go-resty/resty/v2")
	protosPkg  = protogen.GoImportPath("github.com/noble-gase/ne/protos")
	codesPkg   = protogen.GoImportPath("github.com/noble-gase/ne/codes")
)

func protocVersion(p *protogen.Plugin) string {
	v := p.Request.GetCompilerVersion()
	if v == nil {
		return "(unknown)"
	}
	var suffix string
	if s := v.GetSuffix(); s != "" {
		suffix = "-" + s
	}
	return fmt.Sprintf("v%d.%d.%d%s", v.GetMajor(), v.GetMinor(), v.GetPatch(), suffix)
}

// generateFile generates a `xxx_http.pb.go` file containing HTTP service definitions.
func genServiceFile(p *protogen.Plugin, f *protogen.File) *protogen.GeneratedFile {
	if len(f.Services) == 0 {
		return nil
	}
	filename := f.GeneratedFilenamePrefix + suffix
	gf := p.NewGeneratedFile(filename, f.GoImportPath)
	gf.P("// Code generated by protoc-gen-og. DO NOT EDIT.")
	gf.P("// versions:")
	gf.P("// - protoc-gen-og ", version)
	gf.P("// - protoc           ", protocVersion(p))
	if f.Proto.GetOptions().GetDeprecated() {
		gf.P("// ", f.Desc.Path(), " is deprecated.")
	} else {
		gf.P("// source: ", f.Desc.Path())
	}
	gf.P()
	gf.P("package ", f.GoPackageName)
	gf.P()
	genFileContent(f, gf)
	return gf
}

// generateFileContent generates the HTTP service definitions, excluding the package statement.
func genFileContent(f *protogen.File, gf *protogen.GeneratedFile) {
	for _, service := range f.Services {
		genService(gf, service)
	}
}

func genService(gf *protogen.GeneratedFile, service *protogen.Service) {
	serverType := service.GoName + "HttpServer"
	// Server interface
	genServerInterface(gf, service, serverType)
	gf.P()
	// Unimplemented HttpServer
	genServerUnimplement(gf, service, serverType)
	gf.P()
	// Register HttpServer
	genServerRegister(gf, service, serverType)
	// Register HttpServer methods
	genServerMethods(gf, service, serverType)
	gf.P()
	gf.P("// --------------------------------------------- http client ---------------------------------------------")
	gf.P()
	// Client interface
	clientType := service.GoName + "HttpClient"
	// Client interface
	genClientInterface(gf, service, clientType)
	gf.P()
	// New HttpClient
	genClientNew(gf, service, clientType)
	// Register HttpClient methods
	genClientMethods(gf, service, clientType)
	gf.P()
}

// Method(ctx context.Context, in *MethodReq) (*MethodResp, error)
func serverSignature(gf *protogen.GeneratedFile, method *protogen.Method) string {
	var reqArgs []string
	// params
	reqArgs = append(reqArgs, "ctx "+gf.QualifiedGoIdent(ctxPkg.Ident("Context")))
	reqArgs = append(reqArgs, "in *"+gf.QualifiedGoIdent(method.Input.GoIdent))
	// return
	resp := "(*" + gf.QualifiedGoIdent(method.Output.GoIdent) + ", error)"
	return method.GoName + "(" + strings.Join(reqArgs, ", ") + ") " + resp
}

func genServerInterface(gf *protogen.GeneratedFile, service *protogen.Service, serviceType string) {
	gf.P("// ", serviceType, " is the server API definition for ", service.GoName)
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		gf.P("//")
	}
	gf.AnnotateSymbol(serviceType, protogen.Annotation{Location: service.Location})
	// type XXXHttpServer interface {
	gf.P("type ", serviceType, " interface {")
	for _, m := range service.Methods {
		if m.Desc.IsStreamingClient() || m.Desc.IsStreamingServer() {
			continue
		}
		gf.AnnotateSymbol(serviceType+"."+m.GoName, protogen.Annotation{Location: m.Location})
		if m.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
		}
		gf.P(m.Comments.Leading, serverSignature(gf, m))
	}
	gf.P("}")
}

func genServerUnimplement(gf *protogen.GeneratedFile, service *protogen.Service, serviceType string) {
	gf.P("// Unimplemented", serviceType, " should be embedded to have")
	gf.P("// forward compatible implementations.")
	gf.P("//")
	gf.P("// NOTE: this should be embedded by value instead of pointer to avoid a nil")
	gf.P("// pointer dereference when methods are called.")
	gf.P("type Unimplemented", serviceType, " struct{}")
	for _, m := range service.Methods {
		gf.P()
		gf.P("func (Unimplemented", serviceType, ") ", m.GoName, "(context.Context, *", m.Input.GoIdent, ") (*", m.Output.GoIdent, ", error) {")
		gf.P(`return nil, `, errorPkg.Ident("New"), `("method `, m.GoName, ` not implemented")`)
		gf.P("}")
	}
}

func genServerRegister(gf *protogen.GeneratedFile, service *protogen.Service, serviceType string) {
	gf.P("func Register", serviceType, "(r ", chiPkg.Ident("Router"), ", svc ", serviceType, ") {")
	for _, m := range service.Methods {
		rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if ok && rule != nil {
			method, path := getHttpRouter(rule)
			gf.P(strings.TrimSuffix(m.Comments.Leading.String(), "\n"))
			gf.P("r.", method, `("`, path, `", _`, service.GoName, "_", m.GoName, `(svc))`)
			// additional bindings
			for _, bind := range rule.GetAdditionalBindings() {
				method, path := getHttpRouter(bind)
				gf.P("r.", method, `("`, path, `", _`, service.GoName, "_", m.GoName, `(svc))`)
			}
		}
	}
	gf.P("}")
}

func genServerMethods(gf *protogen.GeneratedFile, service *protogen.Service, serviceType string) {
	for _, m := range service.Methods {
		gf.P()
		gf.P(strings.TrimSuffix(m.Comments.Leading.String(), "\n"))
		gf.P("func _", service.GoName, "_", m.GoName, "(svc ", serviceType, ") http.HandlerFunc {")
		gf.P("return func(w ", httpPkg.Ident("ResponseWriter"), ", r *", httpPkg.Ident("Request"), ") {")
		gf.P("ctx := r.Context()")
		gf.P("// parse request")
		gf.P("req := new(", m.Input.GoIdent, ")")
		gf.P("if err := ", contribPkg.Ident("BindProto"), "(r, req); err != nil {")
		gf.P(resultPkg.Ident("Err"), `(ErrParams.New(err.Error())).JSON(w, r)`)
		gf.P("return")
		gf.P("}")
		gf.P("// call service")
		gf.P("resp, err := svc.", m.GoName, "(ctx, req)")
		gf.P("if err != nil {")
		gf.P(resultPkg.Ident("Err"), "(err).JSON(w, r)")
		gf.P("return")
		gf.P("}")
		gf.P(resultPkg.Ident("OK"), "(resp).JSON(w, r)")
		gf.P("}")
		gf.P("}")
	}
}

// Method(ctx context.Context, in *MethodReq, opts ...protos.RequestOption) (*MethodResp, error)
func clientSignature(gf *protogen.GeneratedFile, method *protogen.Method) string {
	var reqArgs []string
	// params
	reqArgs = append(reqArgs, "ctx "+gf.QualifiedGoIdent(ctxPkg.Ident("Context")))
	reqArgs = append(reqArgs, "in *"+gf.QualifiedGoIdent(method.Input.GoIdent))
	reqArgs = append(reqArgs, "opts ..."+gf.QualifiedGoIdent(protosPkg.Ident("RequestOption")))
	// return
	resp := "(*" + gf.QualifiedGoIdent(method.Output.GoIdent) + ", error)"
	return method.GoName + "(" + strings.Join(reqArgs, ", ") + ") " + resp
}

func genClientInterface(gf *protogen.GeneratedFile, service *protogen.Service, serviceType string) {
	gf.P("// ", serviceType, " is the client API definition for ", service.GoName)
	if service.Desc.Options().(*descriptorpb.ServiceOptions).GetDeprecated() {
		gf.P("//")
	}
	gf.AnnotateSymbol(serviceType, protogen.Annotation{Location: service.Location})
	// type XXXHttpClient interface {
	gf.P("type ", serviceType, " interface {")
	for _, m := range service.Methods {
		if m.Desc.IsStreamingClient() || m.Desc.IsStreamingServer() {
			continue
		}
		gf.AnnotateSymbol(serviceType+"."+m.GoName, protogen.Annotation{Location: m.Location})
		if m.Desc.Options().(*descriptorpb.MethodOptions).GetDeprecated() {
		}
		gf.P(m.Comments.Leading, clientSignature(gf, m))
	}
	gf.P("}")
}

func genClientNew(gf *protogen.GeneratedFile, _ *protogen.Service, serviceType string) {
	gf.P("type ", unexport(serviceType), " struct {")
	gf.P("client *", restyPkg.Ident("Client"))
	gf.P("}")
	gf.P()
	gf.P("// New", serviceType, " returns a client for ", serviceType, ". Typically requires WithBaseURL")
	gf.P("func New", serviceType, "(hc *", httpPkg.Ident("Client"), ", opts ...", protosPkg.Ident("ClientOption"), ") ", serviceType, " {")
	gf.P("c := ", restyPkg.Ident("NewWithClient"), "(hc)")
	gf.P("for _, f := range opts {")
	gf.P("f(c)")
	gf.P("}")
	gf.P("return &", unexport(serviceType), "{client: c}")
	gf.P("}")
}

func genClientMethods(gf *protogen.GeneratedFile, service *protogen.Service, serviceType string) {
	for _, m := range service.Methods {
		rule, ok := proto.GetExtension(m.Desc.Options(), annotations.E_Http).(*annotations.HttpRule)
		if !ok || rule == nil {
			continue
		}

		method, path := getHttpRouter(rule)
		isGetMethod := false
		if strings.ToUpper(method) == http.MethodGet {
			isGetMethod = true
		}

		gf.P()
		gf.P(strings.TrimSuffix(m.Comments.Leading.String(), "\n"))
		gf.P("func (c *", unexport(serviceType), ") ", m.GoName, "(ctx ", ctxPkg.Ident("Context"), ", in *", gf.QualifiedGoIdent(m.Input.GoIdent), ", opts ...", protosPkg.Ident("RequestOption"), ") (*"+gf.QualifiedGoIdent(m.Output.GoIdent)+", error) {")
		if isGetMethod {
			gf.P("req := c.client.R().SetContext(ctx).SetQueryParamsFromValues(", protosPkg.Ident("MessageToValues"), "(in))")
		} else {
			gf.P("req := c.client.R().SetContext(ctx)")
		}
		gf.P("for _, f := range opts {")
		gf.P("f(req)")
		gf.P("}")
		if !isGetMethod {
			gf.P("// set request body")
			gf.P("switch ", contribPkg.Ident("ContentType"), "(req.Header) {")
			gf.P("case ", contribPkg.Ident("ContentForm"), ",", contribPkg.Ident("ContentMultipartForm"), ":")
			gf.P("req.SetFormDataFromValues(", protosPkg.Ident("MessageToValues"), "(in))")
			gf.P("default:")
			gf.P("req.SetHeader(", contribPkg.Ident("HeaderContentType"), ", ", contribPkg.Ident("ContentJSON"), ").SetBody(in)")
			gf.P("}")
		}
		gf.P("// send request")
		gf.P("ret := new(", protosPkg.Ident("ApiResult[*"), gf.QualifiedGoIdent(m.Output.GoIdent), "])")
		gf.P("if _, err := req.SetResult(ret).", method, `("`, path, `"); err != nil {`)
		gf.P("return nil, err")
		gf.P("}")
		gf.P("if ret.Code != 0 {")
		gf.P("return nil, ", codesPkg.Ident("New"), "(ret.Code, ret.Msg)")
		gf.P("}")
		gf.P("return ret.Data, nil")
		gf.P("}")
	}
}

func getHttpRouter(rule *annotations.HttpRule) (string, string) {
	switch v := rule.Pattern.(type) {
	case *annotations.HttpRule_Get:
		return "Get", v.Get
	case *annotations.HttpRule_Put:
		return "Put", v.Put
	case *annotations.HttpRule_Post:
		return "Post", v.Post
	case *annotations.HttpRule_Delete:
		return "Delete", v.Delete
	case *annotations.HttpRule_Patch:
		return "Patch", v.Patch
	case *annotations.HttpRule_Custom:
		return v.Custom.GetKind(), v.Custom.GetPath()
	}
	return "Unknown", ""
}

// genCodeFile generates a `code.pb.go` file containing HTTP service definitions.
func genCodeFile(p *protogen.Plugin, f *protogen.File) *protogen.GeneratedFile {
	if !strings.HasSuffix(f.Desc.Path(), "code.proto") || len(f.Enums) == 0 {
		return nil
	}
	filename := f.GeneratedFilenamePrefix + suffix
	gf := p.NewGeneratedFile(filename, f.GoImportPath)
	gf.P("// Code generated by protoc-gen-og. DO NOT EDIT.")
	gf.P("// versions:")
	gf.P("// - protoc-gen-og ", version)
	gf.P("// - protoc           ", protocVersion(p))
	if f.Proto.GetOptions().GetDeprecated() {
		gf.P("// ", f.Desc.Path(), " is deprecated.")
	} else {
		gf.P("// source: ", f.Desc.Path())
	}
	gf.P()
	gf.P("package ", f.GoPackageName)
	gf.P()
	genCodeContent(f, gf)
	return gf
}

// genCodeContent generates the HTTP code definitions, excluding the package statement.
func genCodeContent(f *protogen.File, gf *protogen.GeneratedFile) {
	gf.P("var (")
	for _, e := range f.Enums {
		for _, v := range e.Values {
			msg := strings.ToLower(string(v.Desc.Name()))
			if comment := string(v.Comments.Trailing); len(comment) != 0 {
				msg = strings.TrimSpace(comment)
			}
			name := case2camel(string(v.Desc.Name()))
			gf.P(name, " = ", codesPkg.Ident("New"), "(int(Code_", v.Desc.Name(), `), "`, msg, `")`)
		}
		gf.P()
	}
	gf.P(")")
	gf.P()
	for _, e := range f.Enums {
		for _, v := range e.Values {
			name := case2camel(string(v.Desc.Name()))
			gf.P("func Is", name, "(err error) bool {")
			gf.P("return ", codesPkg.Ident("Is"), "(err, ", name, ")")
			gf.P("}")
			gf.P()
		}
	}
}

func case2camel(s string) string {
	s = strings.ToLower(s)
	items := strings.Split(s, "_")
	count := len(items)
	if count == 1 {
		return export(s)
	}
	for i := 0; i < count; i++ {
		items[i] = export(items[i])
	}
	return strings.Join(items, "")
}

func export(s string) string { return strings.ToUpper(s[:1]) + s[1:] }

func unexport(s string) string { return strings.ToLower(s[:1]) + s[1:] }
