# kmstatus

VERSION := 0.1
LDFLAGS := -s -w -X main.version=$(VERSION)

.PHONY: build run buildx runx test clean install uninstall doc

build:
	@CGO_ENABLED=0 go build -ldflags="$(LDFLAGS)" .

run: build
	@./kmstatus $(ARGS)

buildx:
	@go build -tags X -ldflags="$(LDFLAGS)" .

runx: buildx
	@./kmstatus $(ARGS)

test:
	BASE_PATH="${shell pwd}" go test -cover ./...

clean:
	rm -f kmstatus

# Does not build, so that a binary built with buildx is not replaced.
install:
	@test -f kmstatus || { echo "kmstatus binary not found, run: make build or make buildx"; exit 1; }
	mkdir -p /usr/local/bin
	cp -f kmstatus /usr/local/bin
	chmod 755 /usr/local/bin/kmstatus
	mkdir -p /usr/local/man/man1
	cp kmstatus.1 /usr/local/man/man1/
	chmod 644 /usr/local/man/man1/kmstatus.1

uninstall:
	rm -f /usr/local/bin/kmstatus /usr/local/man/man1/kmstatus.1

# doc.md is the source of the documentation.
# kmstatus.1 (man page) and doc.txt (embedded, printed by --doc) are generated from it.
doc:
	pandoc --standalone --to man \
		-V title=KMSTATUS -V section=1 -V header=DOCUMENTATION \
		-V footer=$(VERSION) -V date="$(shell date +%F)" \
		doc.md -o kmstatus.1
	MANWIDTH=80 man ./kmstatus.1 | col -bx > doc.txt
