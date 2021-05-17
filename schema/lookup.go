package schema

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
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
			m = append(m, search(path))
		}
		path = strings.Replace(m[0], "doc/", "", 1)
		path = strings.Replace(path, "api/", "", 1)
		path = strings.Replace(path, "documentation/", "", 1)
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

func search(s string) string {
	ctx, cancel := chromedp.NewExecAllocator(context.Background(), append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/90.0.4430.212 Safari/537.36"))...)
	defer cancel()
	ctx, cancel = chromedp.NewContext(ctx) //chromedp.WithDebugf(log.Printf)
	defer cancel()
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var nodes []*cdp.Node
	mustRun(ctx,
		chromedp.Navigate("https://developer.apple.com/search/?q="+s),
		chromedp.WaitVisible(`.results-summary`),
		chromedp.Nodes(`.search-result .result-title a`, &nodes),
	)
	return strings.Trim(nodes[0].AttributeValue("href"), "/")
}
