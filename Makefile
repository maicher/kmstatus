# stbar - status bar

build:
	go build .

run: build
	./stbar

test:
	go test -cover ./...

clean:
	rm -f stbar

install: build
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f stbar ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/stbar

uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/stbar
