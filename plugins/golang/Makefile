.PHONY: lint test vendor clean

export GO111MODULE=on

default: lint test

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
