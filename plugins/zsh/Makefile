.PHONY: lint test vendor clean

export GO111MODULE=on

default: lint test

generate_golden_master:
	go test -golden-master -test.v templates_test.go

lint:
	golangci-lint run

test:
	go test -v -cover ./...

yaegi_test:
	yaegi test -v .

vendor:
	go mod vendor
	rm -rdf _gopath && mkdir -p _gopath && cd _gopath && ln -s ../vendor src && cd ..

clean:
	rm -rf ./vendor
