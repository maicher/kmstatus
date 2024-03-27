# kmst

build:
	@CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=0.1" .

run: build
	@./kmst $(ARGS)

buildx:
	@go build -tags X -ldflags="-s -w -X main.version=0.1" .

runx: buildx
	@./kmst $(ARGS)

test:
	BASE_PATH="${shell pwd}" go test -cover ./...

clean:
	rm -f kmst

install:
	mkdir -p /usr/local/bin
	cp -f kmst /usr/local/bin
	chmod 755 /usr/local/bin/kmst
	mkdir -p /usr/local/man/man1
	cp kmst.1 /usr/local/man/man1/
	chmod 644 /usr/local/man/man1/kmst.1

uninstall:
	rm -f /usr/local/bin/kmst /usr/local/man/man1/kmst.1

generatedoc:
	cmd/doc
