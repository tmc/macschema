package schema

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type Lookup struct {
	Query   string
	Lang    string
	Name    string
	Prefix  string
	DocPath string
	APIPath string
	URL     string
}

func (l Lookup) DocExists() bool {
	if _, err := os.Stat(l.DocPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func (l Lookup) APIExists() bool {
	if _, err := os.Stat(l.APIPath); os.IsNotExist(err) {
		return false
	}
	return true
}

func NewLookup(query, lang string) Lookup {
	l := Lookup{
		Query: query,
		Lang:  lang,
	}
	path := strings.ToLower(query)
	if !strings.Contains(path, "/") {
		m, err := filepath.Glob(fmt.Sprintf("./doc/*/%s.%s.json", path, lang))
		fatal(err)
		if len(m) == 0 {
			m, err = filepath.Glob(fmt.Sprintf("./api/*/%s.%s.json", path, lang))
			fatal(err)
		}
		if len(m) == 0 {
			log.Fatal("TODO: search")
		}
		path = strings.Replace(m[0], "doc/", "", 1)
		path = strings.Replace(path, "api/", "", 1)
		path = strings.Replace(path, ".objc.json", "", 1)
		path = strings.Replace(path, ".swift.json", "", 1)
	}
	l.Prefix = filepath.Dir(path)
	l.Name = filepath.Base(path)
	ext := fmt.Sprintf(".%s.json", lang)
	l.DocPath = filepath.Join("./doc", l.Prefix, l.Name+ext)
	l.APIPath = filepath.Join("./api", l.Prefix, l.Name+ext)
	l.URL = fmt.Sprintf("%s%s/%s?language=%s", BaseURL, l.Prefix, l.Name, l.Lang)
	return l
}
