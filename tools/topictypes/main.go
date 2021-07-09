package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/progrium/macschema/schema"
)

func main() {
	var paths []string
	for _, p := range []string{"./doc/*.objc.json", "./doc/*/*.objc.json", "./doc/*/*/*.objc.json"} {
		matches, err := filepath.Glob(p)
		if err != nil {
			log.Fatal(err)
		}
		paths = append(paths, matches...)
	}
	for _, p := range paths {
		var b []byte
		b, err := ioutil.ReadFile(p)
		if err != nil {
			log.Fatal(err)
		}
		var t schema.Topic
		err = json.Unmarshal(b, &t)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(t.Type)
	}
}
