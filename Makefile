# stbar - status bar

build:
	go build .

run:
	go run .

test:
	go test -cover ./...

clean:
	rm -f stbar

install:
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f stbar ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/stbar

uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/stbar
