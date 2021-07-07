package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"

	"github.com/seanpfeifer/rigging/logging"
	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	outFileName = "generated.js"
	prefix      = "rpc/"
)

func main() {
	// Read the code gen request from Stdin, where protoc is writing
	req, err := readRequest(os.Stdin)
	logging.FatalIfError(err)

	// No special options for this generator
	opts := protogen.Options{}
	// Parse out all of the plugin info from the request
	plugin, err := opts.New(req)
	logging.FatalIfError(err)

	// Actually do the generation using the nice structures we get from protogen.Plugin
	// Note that protogen.Plugin already has walked the dependency tree to handle imports
	resp := generatePlugin(plugin)

	// Finally marshal our response to Stdout for the calling protoc to handle
	err = marshalResponse(os.Stdout, resp)
	logging.FatalIfError(err)
}

const twirpUtil = `function createRequest(url, body) {
  return new Request(url, {
    method: "POST",
    credentials: "same-origin",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(body),
  });
}
`

func generatePlugin(plugin *protogen.Plugin) *pluginpb.CodeGeneratorResponse {
	buf := bytes.Buffer{}
	// First write our utility function
	buf.WriteString(twirpUtil)

	for _, file := range plugin.Files {
		for _, svc := range file.Services {
			for _, method := range svc.Methods {
				generatePluginMethod(&buf, method)
			}
		}
	}

	respFile := pluginpb.CodeGeneratorResponse_File{}
	respFile.Content = proto.String(buf.String())
	respFile.Name = proto.String(outFileName)

	resp := new(pluginpb.CodeGeneratorResponse)
	resp.File = []*pluginpb.CodeGeneratorResponse_File{&respFile}
	return resp
}

type jsMethod struct {
	*protogen.Method
	PathPrefix string
}

// JSName exists as a way to get our camelCase method name.
func (j jsMethod) JSName() string {
	if j.GoName == "" {
		return ""
	}
	return strings.ToLower(j.GoName[:1]) + j.GoName[1:]
}

func generatePluginMethod(w io.Writer, method *protogen.Method) {
	funcTempl := `
{{.Comments.Leading}}export async function {{.JSName}}({{range $i, $v := .Input.Fields}}{{if $i}}, {{end}}{{$v.Desc.JSONName}}{{end}}) {
	const res = await fetch(createRequest("/{{.PathPrefix}}{{.Desc.ParentFile.Package}}.{{.Parent.GoName}}/{{.GoName}}", { {{range $i, $v := .Input.Fields}}{{if $i}}, {{end}}"{{$v.Desc.JSONName}}": {{$v.Desc.JSONName}}{{end}} }));
	const jsonBody = await res.json();
	if (res.ok) {
		return jsonBody;
	}
	throw new Error(jsonBody.msg);
}
`

	t := template.Must(template.New("func").Parse(funcTempl))

	in := jsMethod{
		Method:     method,
		PathPrefix: prefix,
	}
	t.Execute(w, in)
}

// readRequest reads the encoded request from the Reader, returning it and an error (if any).
func readRequest(r io.Reader) (*pluginpb.CodeGeneratorRequest, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}

	req := pluginpb.CodeGeneratorRequest{}
	if err = proto.Unmarshal(data, &req); err != nil {
		return nil, fmt.Errorf("error parsing input proto: %v", err)
	}

	if len(req.FileToGenerate) == 0 {
		return nil, errors.New("no files to generate")
	}

	return &req, nil
}

// marshalResponse marshals the resulting response to the given Writer, returning an error if any occurs.
func marshalResponse(w io.Writer, resp *pluginpb.CodeGeneratorResponse) error {
	data, err := proto.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshaling response: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		return fmt.Errorf("writing response: %v", err)
	}

	return nil
}
