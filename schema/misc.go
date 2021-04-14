package schema

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/chromedp/chromedp"
)

const (
	BaseURL = "https://developer.apple.com/documentation/"
	Version = 2
)

func LookupFromPath(path string) Lookup {
	u, err := url.Parse(path)
	if err != nil {
		panic(err)
	}
	query := strings.Replace(u.Path, "/documentation/", "", 1)
	return NewLookup(query, u.Query().Get("language"))
}

func ReadTopic(l Lookup) (t Topic, err error) {
	var b []byte
	b, err = ioutil.ReadFile(l.DocPath)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &t)
	return
}

func Stats() {
	stats := make(map[string]int)
	m, err := filepath.Glob("./documentation/**/**.objc.json")
	fatal(err)
	for _, match := range m {
		b, err := ioutil.ReadFile(match)
		fatal(err)
		var t Topic
		fatal(json.Unmarshal(b, &t))
		stats[t.Type]++
	}
	fmt.Println(stats)
}

func Types(path string) {
	c := readSchema(path)
	var types []DataType
	collectTypes(&types, reflect.ValueOf(c))
	uniq := make(map[string]bool)
	for _, t := range types {
		uniq[t.Name] = true
	}
	for k := range uniq {
		fmt.Println(k)
	}
}

func mustRun(ctx context.Context, actions ...chromedp.Action) {
	if err := chromedp.Run(ctx, actions...); err != nil {
		log.Fatal(err)
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func collectTypes(types *[]DataType, src reflect.Value) {
	if src.Kind() == reflect.Ptr {
		src = src.Elem()
	}
	typeInfo := reflect.TypeOf(DataType{})
	switch src.Kind() {
	case reflect.Struct:
		for i := 0; i < src.NumField(); i += 1 {
			f := src.Field(i)
			if f.Type() == typeInfo {
				*types = append(*types, f.Interface().(DataType))
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

func readSchema(path string) Class {
	b, err := ioutil.ReadFile(path)
	fatal(err)

	var c Class
	fatal(json.Unmarshal(b, &c))
	return c
}
