package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/progrium/macschema/pkg/declparse"
	"github.com/progrium/macschema/pkg/topic"
)

type Class struct {
	Name        string
	Description string
	Declaration string

	InstanceMethods    []Method   `json:",omitempty"`
	InstanceProperties []Property `json:",omitempty"`

	TypeMethods    []Method   `json:",omitempty"`
	TypeProperties []Property `json:",omitempty"`

	Frameworks []string
	Platforms  []string

	Deprecated bool `json:",omitempty"`
	URL        string

	ParseDate    time.Time
	ParseVersion int
}

type TypeInfo struct {
	Name        string     `json:",omitempty"`
	IsPtr       bool       `json:",omitempty"`
	IsPtrPtr    bool       `json:",omitempty"`
	Annotations []string   `json:",omitempty"`
	Func        *Func      `json:",omitempty"`
	Params      []TypeInfo `json:",omitempty"`
}

func TypeInfoFromAst(ti declparse.TypeInfo) TypeInfo {
	var fn *Func
	if ti.Func != nil {
		fn = FuncFromAst(ti.Func)
	}
	var params []TypeInfo
	for _, param := range ti.Params {
		params = append(params, TypeInfoFromAst(param))
	}
	var annots []string
	for annot := range ti.Annots {
		annots = append(annots, strings.ToLower(annot.String()))
	}
	return TypeInfo{
		Name:        ti.Name,
		IsPtr:       ti.IsPtr,
		IsPtrPtr:    ti.IsPtrPtr,
		Annotations: annots,
		Func:        fn,
		Params:      params,
	}
}

type Func struct {
	Name       string `json:",omitempty"`
	ReturnType TypeInfo
	Args       []ArgInfo
}

func FuncFromAst(fn *declparse.FunctionDecl) *Func {
	var args []ArgInfo
	for _, arg := range fn.Args {
		args = append(args, ArgInfoFromAst(arg))
	}
	return &Func{
		Name:       fn.Name,
		ReturnType: TypeInfoFromAst(fn.ReturnType),
		Args:       args,
	}
}

type ArgInfo struct {
	Name string `json:",omitempty"`
	Type TypeInfo
}

func ArgInfoFromAst(ai declparse.ArgInfo) ArgInfo {
	return ArgInfo{
		Name: ai.Name,
		Type: TypeInfoFromAst(ai.Type),
	}
}

type Property struct {
	Name        string
	Description string
	Declaration string
	Type        TypeInfo
	Attrs       map[string]interface{}
	Deprecated  bool `json:",omitempty"`
	URL         string
}

func PropertyFromAst(p declparse.PropertyDecl) Property {
	prop := Property{
		Name: p.Name,
		Type: TypeInfoFromAst(p.Type),
	}
	attrs := make(map[string]interface{})
	for attr, val := range p.Attrs {
		var v interface{}
		if val == "" {
			v = true
		} else {
			v = val
		}
		attrs[attr.String()] = v
	}
	prop.Attrs = attrs
	return prop
}

type Method struct {
	Name        string
	Description string
	Declaration string
	Return      TypeInfo
	Args        []ArgInfo
	Deprecated  bool `json:",omitempty"`
	URL         string
}

func MethodFromAst(m declparse.MethodDecl) Method {
	var args []ArgInfo
	for _, arg := range m.Args {
		args = append(args, ArgInfoFromAst(arg))
	}
	return Method{
		Name:   m.Name(),
		Return: TypeInfoFromAst(m.ReturnType),
		Args:   args,
	}
}

const BaseURL = "https://developer.apple.com"

const Version = 1

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func readTopicFromURL(pathOrUrl string) topic.Topic {
	if strings.HasPrefix(pathOrUrl, "/") {
		pathOrUrl = BaseURL + pathOrUrl
	}
	u, _ := url.Parse(pathOrUrl)
	lang := u.Query().Get("language")
	if lang == "" {
		lang = "swift"
	}
	return readTopic(fmt.Sprintf(".%s.%s.json", u.Path, lang))
}

func collectTypes(types *[]TypeInfo, src reflect.Value) {
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	typeInfo := reflect.TypeOf(TypeInfo{})
	switch src.Kind() {
	case reflect.Struct:
		for i := 0; i < src.NumField(); i += 1 {
			f := src.Field(i)
			if f.Type() == typeInfo {
				*types = append(*types, f.Interface().(TypeInfo))
			} else {
				collectTypes(types, f)
			}
		}
	case reflect.Slice:
		for i := 0; i < src.Len(); i += 1 {
			collectTypes(types, src.Index(i))
		}
	}
}

func readTopic(path string) topic.Topic {
	b, err := ioutil.ReadFile(path)
	fatal(err)

	var t topic.Topic
	fatal(json.Unmarshal(b, &t))
	return t
}

func readSchema(path string) Class {
	b, err := ioutil.ReadFile(path)
	fatal(err)

	var c Class
	fatal(json.Unmarshal(b, &c))
	return c
}

func Types(path string) {
	c := readSchema(path)
	var types []TypeInfo
	collectTypes(&types, reflect.ValueOf(c))
	uniq := make(map[string]bool)
	for _, t := range types {
		uniq[t.Name] = true
	}
	for k := range uniq {
		fmt.Println(k)
	}
}

func Parse(path string) {
	t := readTopic(path)

	var c Class
	c.ParseDate = t.LastFetch
	c.ParseVersion = Version
	c.Declaration = t.Declaration
	c.Name = t.Title
	c.Description = t.Description
	c.Frameworks = t.Frameworks
	c.URL = BaseURL + t.Path
	for _, p := range t.Platforms {
		if p == "Deprecated" {
			c.Deprecated = true
		} else {
			c.Platforms = append(c.Platforms, p)
		}
	}
	for _, topic := range t.Topics {
		t := readTopicFromURL(topic.Path)
		if t.Type == "Function" || t.Type == "Enumeration" || t.Type == "Global Variable" || t.Type == "Enumeration Case" {
			continue
		}
		var isDeprecated bool
		for _, p := range t.Platforms {
			if p == "Deprecated" {
				isDeprecated = true
			}
		}
		if t.Declaration != "" {
			p := declparse.NewStringParser(t.Declaration)
			ast, err := p.Parse()
			if err != nil {
				if strings.Contains(err.Error(), "typedef") ||
					strings.Contains(err.Error(), "const") {
					continue
				}
				fatal(fmt.Errorf("%s: %w", topic.Path, err))
			}
			switch t.Type {
			case "Type Method":
				m := MethodFromAst(*ast.Method)
				m.Description = t.Description
				m.Declaration = t.Declaration
				m.URL = BaseURL + t.Path
				m.Deprecated = isDeprecated
				c.TypeMethods = append(c.TypeMethods, m)
			case "Instance Method":
				m := MethodFromAst(*ast.Method)
				m.Description = t.Description
				m.Declaration = t.Declaration
				m.URL = BaseURL + t.Path
				m.Deprecated = isDeprecated
				c.InstanceMethods = append(c.InstanceMethods, m)
			case "Type Property":
				p := PropertyFromAst(*ast.Property)
				p.Description = t.Description
				p.Declaration = t.Declaration
				p.URL = BaseURL + t.Path
				p.Deprecated = isDeprecated
				c.TypeProperties = append(c.TypeProperties, p)
			case "Instance Property":
				p := PropertyFromAst(*ast.Property)
				p.Description = t.Description
				p.Declaration = t.Declaration
				p.URL = BaseURL + t.Path
				p.Deprecated = isDeprecated
				c.InstanceProperties = append(c.InstanceProperties, p)
			default:
			}
		}
	}

	b, err := json.MarshalIndent(c, "", "  ")
	fatal(err)
	fatal(ioutil.WriteFile(strings.Replace(path, "documentation", "schema", 1), b, 0644))
}
