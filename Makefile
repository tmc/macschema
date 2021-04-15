
build: local/bin/macschema

clean:
	rm -fr dist
	rm -f local/bin/macschema

local/bin/macschema: */*.go
	go build -ldflags $(ldflags) -o ./local/bin/macschema .


version=$(shell cat version)
module=$(shell go list -m)
branch=$(shell git branch --show-current)
ldflags="-X $(module)/cmd.Version=$(version:dev=$(branch))"