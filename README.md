[![Go Report Card](https://goreportcard.com/badge/github.com/stutonk/pen)](https://goreportcard.com/report/github.com/stutonk/pen)  
pen is a lightweight alternatve to PGP for securely encrypting files; it uses
the Argon2 algorithm for key-stretching passwords and NaCl's secretbox for
data encryption. If a filename ends in `.pen`, it will be decrypted rather
than encrypted. Upon successful encryption or decryption, the original file
or .pen file, respectively, is removed.

```
usage: pen [-h, -v] file [files...]
Options are:
  -h, --help      display this help and exit
  -v, --version   output version information and exit
```

### release binaries
are available [here](https://github.com/stutonk/pen/releases) for amd64/all major OSes

### for unixes
`make && make install`

### everybody else
`go build`

### note
Due to the nature of underlying filesystems and their underlying hardware,
no attempt is made to shred the original file before it is removed following
encryption. A determined person or law enforcement agency with access to
your raw disk WILL be able to recover the original file. Therefore, it is
advised that for critical applications you use a more mature solution that
follows stricter security protocols.