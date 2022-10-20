package main

import (
	"bytes"
	_ "embed"
	"flag"
	"strings"
	"text/template"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/reflect/protoreflect"
	pluginpb "google.golang.org/protobuf/types/pluginpb"
)

const (
	outFileName = "generated.ts"
)

//go:embed template.tmpl
var fileTemplate string

//go:embed listfields.tmpl
var listTemplate string

func main() {
	// Set up our flags. The only one we care about for now is the server path prefix.
	var flags flag.FlagSet
	prefix := flags.String("pathPrefix", "/twirp", "the server path prefix to use, if modified from the Twirp default")

	// No special options for this generator
	opts := protogen.Options{ParamFunc: flags.Set}
	opts.Run(func(plugin *protogen.Plugin) error {
		plugin.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		out := plugin.NewGeneratedFile(outFileName, "")

		in := jsData{
			Files:      plugin.Files,
			PathPrefix: *prefix,
		}

		tsTemplate, err := template.New("file").
			Funcs(template.FuncMap{
				"JSName":  JSName,
				"GetType": in.GetType,
			}).
			Parse(listTemplate)
		if err != nil {
			return err
		}

		// Add onto the list template by adding the rest of the file template
		tsTemplate, err = tsTemplate.Parse(fileTemplate)
		if err != nil {
			return err
		}

		return tsTemplate.Execute(out, in)
	})
}

type jsData struct {
	Files      []*protogen.File
	PathPrefix string
}

func (j *jsData) GetType(desc protoreflect.FieldDescriptor) string {
	switch desc.Kind() {
	case protoreflect.BoolKind:
		return "boolean"
	case protoreflect.EnumKind:
		return "number"
	case protoreflect.Int32Kind:
		return "number"
	case protoreflect.Sint32Kind:
		return "number"
	case protoreflect.Uint32Kind:
		return "number"
	case protoreflect.Sfixed32Kind:
		return "number"
	case protoreflect.Fixed32Kind:
		return "number"
	case protoreflect.FloatKind:
		return "number"
	case protoreflect.Int64Kind:
		return "bigint"
	case protoreflect.Sint64Kind:
		return "bigint"
	case protoreflect.Uint64Kind:
		return "bigint"
	case protoreflect.Sfixed64Kind:
		return "bigint"
	case protoreflect.Fixed64Kind:
		return "bigint"
	case protoreflect.DoubleKind:
		return "number"
	case protoreflect.StringKind:
		return "string"
	case protoreflect.BytesKind:
		// NOT using Uint8Array here, as these end up getting encoded/decoded as base64 strings
		return "string"
	case protoreflect.MessageKind:
		buf := &bytes.Buffer{}
		buf.WriteString("{ ")
		fields := desc.Message().Fields()
		for i := 0; i < fields.Len(); i++ {
			field := fields.Get(i)
			if i > 0 {
				buf.WriteString(", ")
			}
			buf.WriteString(field.JSONName())
			buf.WriteString("?: ")
			buf.WriteString(j.GetType(field))
		}
		buf.WriteString(" }")
		return buf.String()
	case protoreflect.GroupKind: // Not supported - explicitly a deprecated Protobuf feature
		return "any"
	default:
		return "any"
	}
}

// JSName exists as a way to get our camelCase method name.
func JSName(m *protogen.Method) string {
	if m.GoName == "" {
		return ""
	}
	return strings.ToLower(m.GoName[:1]) + m.GoName[1:]
}
