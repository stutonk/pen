#!/bin/sh

TMP=$(mktemp)
HTML_IN=template.html
HTML_OUT=index.html
MD_IN=template.md
MD_OUT=README.md


nroff -man ../pen.1 | col -bx > ${TMP}
sed '/CONTENT/ {
    r '${TMP}'
    d
}' <${HTML_IN}>${HTML_OUT}

[ -e ${TMP} ] && rm ${TMP}

mv ${HTML_OUT} ../docs