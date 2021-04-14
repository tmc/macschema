package schema

import (
	"fmt"
	"strings"

	"github.com/progrium/macschema/declparse"
)

func PullSchema(l Lookup) Schema {
	t, err := ReadTopic(l)
	fatal(err)

	var s Schema
	s.PullDate = t.LastFetch
	s.Version = Version
	s.Kind = "class"

	var c Class
	c.Declaration = t.Declaration
	c.Name = t.Title
	c.Description = t.Description
	c.Frameworks = t.Frameworks
	c.TopicURL = BaseURL + strings.Replace(t.Path, "/documentation/", "", 1)
	for _, p := range t.Platforms {
		if p == "Deprecated" {
			c.Deprecated = true
		} else {
			c.Platforms = append(c.Platforms, p)
		}
	}
	for _, topic := range t.Topics {
		t, err := ReadTopic(LookupFromPath(topic.Path))
		fatal(err)
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
			url := BaseURL + strings.Replace(t.Path, "/documentation/", "", 1)
			switch t.Type {
			case "Type Method":
				m := MethodFromAst(*ast.Method)
				m.Description = t.Description
				m.Declaration = t.Declaration
				m.TopicURL = url
				m.Deprecated = isDeprecated
				c.TypeMethods = append(c.TypeMethods, m)
			case "Instance Method":
				m := MethodFromAst(*ast.Method)
				m.Description = t.Description
				m.Declaration = t.Declaration
				m.TopicURL = url
				m.Deprecated = isDeprecated
				c.InstanceMethods = append(c.InstanceMethods, m)
			case "Type Property":
				p := PropertyFromAst(*ast.Property)
				p.Description = t.Description
				p.Declaration = t.Declaration
				p.TopicURL = url
				p.Deprecated = isDeprecated
				c.TypeProperties = append(c.TypeProperties, p)
			case "Instance Property":
				p := PropertyFromAst(*ast.Property)
				p.Description = t.Description
				p.Declaration = t.Declaration
				p.TopicURL = url
				p.Deprecated = isDeprecated
				c.InstanceProperties = append(c.InstanceProperties, p)
			default:
			}
		}
	}

	s.Class = &c
	return s
}
