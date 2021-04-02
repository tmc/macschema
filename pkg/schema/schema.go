package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/progrium/macschema/pkg/declparser"
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
	Name     string
	IsPtr    bool       `json:",omitempty"`
	IsConst  bool       `json:",omitempty"`
	IsKindOf bool       `json:",omitempty"`
	Block    *Block     `json:",omitempty"`
	Params   []TypeInfo `json:",omitempty"`
}

func TypeInfoFromAst(ti declparser.TypeInfo) TypeInfo {
	return TypeInfo{
		Name:     ti.Name,
		IsPtr:    ti.IsPtr,
		IsConst:  ti.IsConst,
		IsKindOf: ti.IsKindOf,
		// TODO: more
	}
}

type Block struct {
	Name       string
	ReturnType TypeInfo
	Args       []ArgInfo
}

type ArgInfo struct {
	Name string
	Type TypeInfo
}

func ArgInfoFromAst(ai declparser.ArgInfo) ArgInfo {
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
	Attrs       struct {
		Class     bool   `json:",omitempty"`
		Readonly  bool   `json:",omitempty"`
		Weak      bool   `json:",omitempty"`
		Nonatomic bool   `json:",omitempty"`
		Copy      bool   `json:",omitempty"`
		Nullable  bool   `json:",omitempty"`
		Nonnull   bool   `json:",omitempty"`
		Retain    bool   `json:",omitempty"`
		Getter    string `json:",omitempty"`
		Setter    string `json:",omitempty"`
	}
	Deprecated bool `json:",omitempty"`
	URL        string
}

func PropertyFromAst(p declparser.PropertyDecl) Property {
	prop := Property{
		Name: p.Name,
		Type: TypeInfoFromAst(p.Type),
	}
	prop.Attrs.Class = p.Class
	prop.Attrs.Copy = p.Copy
	prop.Attrs.Getter = p.Getter
	prop.Attrs.Nonatomic = p.Nonatomic
	prop.Attrs.Nonnull = p.Nonnull
	prop.Attrs.Nullable = p.Nullable
	prop.Attrs.Readonly = p.Readonly
	prop.Attrs.Setter = p.Setter
	prop.Attrs.Weak = p.Weak
	prop.Attrs.Retain = p.Retain
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

func MethodFromAst(m declparser.MethodDecl) Method {
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

func readTopic(path string) topic.Topic {
	b, err := ioutil.ReadFile(path)
	fatal(err)

	var t topic.Topic
	fatal(json.Unmarshal(b, &t))
	return t
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
		if t.Type == "Function" || t.Type == "Enumeration" || t.Type == "Global Variable" {
			continue
		}
		if t.Declaration != "" {
			p := declparser.NewStringParser(t.Declaration)
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
				c.TypeMethods = append(c.TypeMethods, m)
			case "Instance Method":
				m := MethodFromAst(*ast.Method)
				m.Description = t.Description
				m.Declaration = t.Declaration
				m.URL = BaseURL + t.Path
				c.InstanceMethods = append(c.InstanceMethods, m)
			case "Type Property":
				p := PropertyFromAst(*ast.Property)
				p.Description = t.Description
				p.Declaration = t.Declaration
				p.URL = BaseURL + t.Path
				c.TypeProperties = append(c.TypeProperties, p)
			case "Instance Property":
				p := PropertyFromAst(*ast.Property)
				p.Description = t.Description
				p.Declaration = t.Declaration
				p.URL = BaseURL + t.Path
				c.InstanceProperties = append(c.InstanceProperties, p)
			default:
			}
		}
	}

	b, err := json.MarshalIndent(c, "", "  ")
	fatal(err)
	fmt.Println(string(b))
}
