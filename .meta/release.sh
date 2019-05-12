#!/bin/sh

ARCH=amd64
BIN=pen

for OS in linux darwin freebsd
do
    GOARCH=${ARCH} GOOS=${OS} go build
    tar czf ${BIN}-${ARCH}-${OS}.tar.gz ./${BIN}
done

GOARCH=${ARCH} GOOS=windows go build
zip ${BIN}-${ARCH}-windows.zip ./${BIN}.exe

test -e ./${BIN} && rm -f ./${BIN}
test -e ./${BIN}.exe && rm -f ./${BIN}.exe
