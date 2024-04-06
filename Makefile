# kmstatus

build:
	@CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=0.1" .

run: build
	@./kmstatus $(ARGS)

buildx:
	@go build -tags X -ldflags="-s -w -X main.version=0.1" .

runx: buildx
	@./kmstatus $(ARGS)

test:
	BASE_PATH="${shell pwd}" go test -cover ./...

clean:
	rm -f kmstatus

install:
	mkdir -p /usr/local/bin
	cp -f kmstatus /usr/local/bin
	chmod 755 /usr/local/bin/kmstatus
	mkdir -p /usr/local/man/man1
	cp kmstatus.1 /usr/local/man/man1/
	chmod 644 /usr/local/man/man1/kmstatus.1

uninstall:
	rm -f /usr/local/bin/kmstatus /usr/local/man/man1/kmstatus.1

generatedoc:
	cmd/doc
