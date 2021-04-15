
build: local/bin/macschema

clean:
	rm local/bin/macschema

ldflags="-X $(module)/cmd.Version=$(version:dev=$(branch))"
local/bin/macschema: */*.go
	go build -ldflags $(ldflags) -o ./local/bin/macschema .


version=$(shell cat version)
_module=$(shell head -1 go.mod)
module=$(_module:module=)
branch=$(shell git branch --show-current)