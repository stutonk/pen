PREFIX ?= /usr/local
BINDIR ?= ${PREFIX}/bin
MANDIR ?= ${PREFIX}/share/man

pen:
	go build

install: pen
	install -d ${DESTDIR}${BINDIR}
	install -m 755 pen ${DESTDIR}${BINDIR}
	install -d ${DESTDIR}${MANDIR}/man1
	install -m 644 pen.1 ${DESTDIR}${MANDIR}/man1

uninstall:
	rm ${DESTDIR}${BINDIR}/pen
	rm ${DESTDIR}${MANDIR}/man1/pen.1

clean:
	rm -f pen
