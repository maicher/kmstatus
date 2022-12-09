# kmstatus

build:
	go build .

run: build
	./kmstatus

test:
	go test -cover ./...

clean:
	rm -f kmstatus

install: build
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f kmstatus ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/kmstatus

uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/kmstatus
