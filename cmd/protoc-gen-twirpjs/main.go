package main

import (
	_ "embed"
	"flag"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	pluginpb "google.golang.org/protobuf/types/pluginpb"
)

const (
	outFileName = "generated.js"
)

//go:embed template.tmpl
var fileTemplate string

func main() {
	// Set up our flags. The only one we care about for now is the server path prefix.
	var flags flag.FlagSet
	prefix := flags.String("pathPrefix", "/twirp", "the server path prefix to use, if modified from the Twirp default")

	// No special options for this generator
	opts := protogen.Options{ParamFunc: flags.Set}
	opts.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		out := plugin.NewGeneratedFile(outFileName, "")

		template, err := template.New("file").
			Funcs(template.FuncMap{"JSName": JSName}).
			Parse(fileTemplate)
		if err != nil {
			return err
		}

		in := jsData{
			Files:      plugin.Files,
			PathPrefix: *prefix,
		}

		return template.Execute(out, in)
	})
}

type jsData struct {
	Files      []*protogen.File
	PathPrefix string
}

// JSName exists as a way to get our camelCase method name.
func JSName(m *protogen.Method) string {
	if m.GoName == "" {
		return ""
	}
	return strings.ToLower(m.GoName[:1]) + m.GoName[1:]
}
