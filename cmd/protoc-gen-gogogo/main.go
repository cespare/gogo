package main

import (
	"strings"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	files := req.GetProtoFile()
	files = vanity.FilterFiles(files, vanity.NotInPackageGoogleProtobuf)

	vanity.ForEachFile(files, vanity.TurnOnMarshalerAll)
	vanity.ForEachFile(files, vanity.TurnOnSizerAll)
	vanity.ForEachFile(files, vanity.TurnOnUnmarshalerAll)

	vanity.ForEachFieldInFiles(files, FixFieldNames)

	resp := command.Generate(req)
	command.Write(resp)
}

func FixFieldNames(field *descriptor.FieldDescriptorProto) {
	if FieldHasStringExtension(field, gogoproto.E_Customname) {
		return
	}
	name, ok := SnakeToCamel(field.GetName())
	if !ok {
		return
	}
	if field.Options == nil {
		field.Options = &descriptor.FieldOptions{}
	}
	if err := proto.SetExtension(field.Options, gogoproto.E_Customname, &name); err != nil {
		panic(err)
	}
}

func FieldHasStringExtension(field *descriptor.FieldDescriptorProto, extension *proto.ExtensionDesc) bool {
	if field.Options == nil {
		return false
	}
	value, err := proto.GetExtension(field.Options, extension)
	if err != nil {
		return false
	}
	if value == nil {
		return false
	}
	if value.(*string) == nil {
		return false
	}
	return true
}

func SnakeToCamel(snake string) (camel string, changed bool) {
	parts := strings.Split(snake, "_")
	for _, part := range parts {
		if commonInitialisms[part] {
			changed = true
			camel += strings.ToTitle(part)
		} else {
			camel += strings.Title(part)
		}
	}
	return camel, changed
}

// Taken from golint (with some additions).
var commonInitialisms = map[string]bool{
	"api":   true,
	"ascii": true,
	"cpu":   true,
	"css":   true,
	"dns":   true,
	"eof":   true,
	"guid":  true,
	"html":  true,
	"http":  true,
	"https": true,
	"id":    true,
	"ip":    true,
	"json":  true,
	"lhs":   true,
	"os":    true,
	"qps":   true,
	"ram":   true,
	"rhs":   true,
	"rpc":   true,
	"sla":   true,
	"smtp":  true,
	"sql":   true,
	"ssh":   true,
	"tcp":   true,
	"tls":   true,
	"ttl":   true,
	"udp":   true,
	"ui":    true,
	"uid":   true,
	"uri":   true,
	"url":   true,
	"utf8":  true,
	"uuid":  true,
	"vm":    true,
	"xml":   true,
	"xsrf":  true,
	"xss":   true,
}
