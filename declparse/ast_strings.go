package declparse

import (
	"fmt"
	"strings"
)

func (s Statement) String() string {
	if s.Method != nil {
		return s.Method.String() + ";"
	}
	if s.Property != nil {
		return s.Property.String() + ";"
	}
	if s.Interface != nil {
		return s.Interface.String() + ";"
	}
	if s.Protocol != nil {
		return s.Protocol.String() + ";"
	}
	if s.Function != nil {
		return s.Function.String() + ";"
	}
	if s.Variable != nil {
		return s.Variable.String() + ";"
	}
	if s.Enum != nil {
		return s.Enum.String() + ";"
	}
	return ""
}

func (i ProtocolDecl) String() string {
	b := &strings.Builder{}
	_, _ = fmt.Fprintf(b, "@protocol %s", i.Name)
	if i.SuperName != "" {
		_, _ = fmt.Fprintf(b, " : %s", i.SuperName)
	}
	return b.String()
}

func (i InterfaceDecl) String() string {
	b := &strings.Builder{}
	_, _ = fmt.Fprintf(b, "@interface %s", i.Name)
	if i.SuperName != "" {
		_, _ = fmt.Fprintf(b, " : %s", i.SuperName)
	}
	return b.String()
}

func (p PropertyDecl) String() string {
	b := &strings.Builder{}
	b.WriteString("@property")
	var attrs []string
	for _, attr := range PropAttrs() {
		val, ok := p.Attrs[attr]
		if !ok {
			continue
		}
		if val == "" {
			attrs = append(attrs, attr.String())
		} else {
			attrs = append(attrs, fmt.Sprintf("%s=%s", attr, val))
		}
	}
	if len(attrs) != 0 {
		fmt.Fprintf(b, "(%s)", strings.Join(attrs, ", "))
	}
	b.WriteString(" ")
	typ := p.Type.String()
	b.WriteString(typ)
	if typ[len(typ)-1] != '*' {
		b.WriteString(" ")
	}
	b.WriteString(p.Name)
	return b.String()
}

func (args FuncArgs) String() string {
	var str []string
	for _, arg := range args {
		str = append(str, strings.Trim(fmt.Sprintf("%s %s", arg.Type, arg.Name), " "))
	}
	return strings.Join(str, ", ")
}

func (f FunctionDecl) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%s %s(%s)", f.ReturnType, f.Ident(), f.Args)
	return b.String()
}

func (m MethodDecl) String() string {
	b := &strings.Builder{}
	prefix := "-"
	if m.TypeMethod {
		prefix = "+"
	}
	if len(m.Args) == 0 {
		fmt.Fprintf(b, "%s (%s)%s", prefix, m.ReturnType, m.Name())
	} else {
		var parts []string
		for arg, part := range m.NameParts {
			parts = append(parts, fmt.Sprintf("%s:%s", part, m.Args[arg]))
		}
		fmt.Fprintf(b, "%s (%s)%s", prefix, m.ReturnType, strings.Join(parts, " "))
		if m.Args[len(m.Args)-1].Name == "..." {
			b.WriteString(", ...")
		}
	}
	return b.String()
}

func (t TypeInfo) String() string {
	if t.Func != nil {
		return t.Func.String()
	}
	params := ""
	if len(t.Params) > 0 {
		var p []string
		for _, param := range t.Params {
			p = append(p, param.String())
		}
		params = fmt.Sprintf("<%s>", strings.Join(p, ", "))
	}
	ptr := ""
	if t.IsPtr {
		ptr = "*"
	}
	str := strings.Trim(fmt.Sprintf("%s%s %s", t.Name, params, ptr), " ")
	for annot, ok := range t.Annots {
		if !ok {
			continue
		}
		str = fmt.Sprintf(annot.Format(), str)
	}
	b := &strings.Builder{}
	b.WriteString(str)
	if t.IsPtrPtr {
		if str[len(str)-1] != '*' {
			b.WriteString(" ")
		}
		b.WriteString("*")
	}
	return b.String()
}

func (arg ArgInfo) String() string {
	return fmt.Sprintf("(%s)%s", arg.Type, arg.Name)
}

func (v VariableDecl) String() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%s %s", v.Type, v.Name)
	if v.Value != "" {
		fmt.Fprintf(b, " = %s", v.Value)
	}
	return b.String()
}

func (e EnumDecl) String() string {
	b := &strings.Builder{}
	if e.Name != "" {
		fmt.Fprintf(b, "enum %s { ", e.Name)
	} else {
		b.WriteString("enum { ")
	}
	for idx, c := range e.Consts {
		if c.Value != "" {
			fmt.Fprintf(b, "%s = %s", c.Name, c.Value)
		} else {
			fmt.Fprintf(b, "%s", c.Name)
		}
		if idx == len(e.Consts)-1 {
			b.WriteString(" ")
		} else {
			b.WriteString(", ")
		}
	}
	b.WriteString("}")
	return b.String()
}
