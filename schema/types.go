package schema

import (
	"time"
)

type Schema struct {
	Class *Class `json:",omitempty"`

	Kind     string
	PullDate time.Time
	Version  int
}

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

	Deprecated bool   `json:",omitempty"`
	TopicURL   string `json:",omitempty"`
}

type DataType struct {
	Name        string     `json:",omitempty"`
	IsPtr       bool       `json:",omitempty"`
	IsPtrPtr    bool       `json:",omitempty"`
	Annotations []string   `json:",omitempty"`
	FuncPtr     *Func      `json:",omitempty"`
	Block       *Func      `json:",omitempty"`
	Params      []DataType `json:",omitempty"`
}

type Func struct {
	Name     string `json:",omitempty"`
	Return   DataType
	Args     []Arg
	TopicURL string `json:",omitempty"`
}

type Arg struct {
	Name string `json:",omitempty"`
	Type DataType
}

type Property struct {
	Name        string
	Description string
	Declaration string
	Type        DataType
	Attrs       map[string]interface{}
	Deprecated  bool   `json:",omitempty"`
	TopicURL    string `json:",omitempty"`
}

type Method struct {
	Name        string
	Description string
	Declaration string
	Return      DataType
	Args        []Arg
	Deprecated  bool   `json:",omitempty"`
	TopicURL    string `json:",omitempty"`
}

type Topic struct {
	Path        string
	Title       string
	Type        string
	Description string
	Declaration string
	Frameworks  []string
	Platforms   []string
	Topics      []Link
	LastFetch   time.Time
	LastVersion int
}

type Link struct {
	Section string
	Name    string
	Path    string
}
