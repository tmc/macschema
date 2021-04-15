
build: local/bin/macschema

clean:
	rm -fr dist
	rm -f local/bin/macschema

release:
	git tag v$(version:dev=)
	git push origin v$(version:dev=)
	goreleaser release --rm-dist
	@echo "==> Remember to update ./version! Current contents: $(version)"
	

local/bin/macschema: */*.go
	go build -ldflags $(ldflags) -o ./local/bin/macschema .


version=$(shell cat version)
module=$(shell go list -m)
branch=$(shell git branch --show-current)
ldflags="-X $(module)/cmd.Version=$(version:dev=$(branch))"