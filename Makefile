# kmstatus

build:
	go build .

run: build
	./kmstatus

runx:
	go run -tags X . -x

buildx:
	go build -tags X .

test:
	go test -cover ./...

clean:
	rm -f kmstatus

install:
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f kmstatus ${DESTDIR}${PREFIX}/bin
	chmod 755 ${DESTDIR}${PREFIX}/bin/kmstatus

uninstall:
	rm -f ${DESTDIR}${PREFIX}/bin/kmstatus
