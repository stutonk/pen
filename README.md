[![Go Report Card](https://goreportcard.com/badge/github.com/stutonk/pen)](https://goreportcard.com/report/github.com/stutonk/pen)  
pen is a lightweight alternatve to PGP for securely encrypting files; it uses
the Argon2 algorithm for key-stretching passwords and NaCl's secretbox for
data encryption. If a filename ends in `.pen`, it will be decrypted rather
than encrypted. Upon successful encryption or decryption, the original file
or .pen file, respectively, is securely removed using shred from GNU coreutils.
If shred is not available, the original file remains untouched.
```
usage: pen [-h, -v] file [files...]
Options are:
  -h, --help      display this help and exit
  -v, --version   output version information and exit
```

### release binaries
are available [here](https://github.com/stutonk/pen/releases) for amd64/all
major OSes

### for unixes
`make && make install`

### everybody else
`go build`

### note
Shredding the original input may still be insecure due to the underlying
charactersitics of your filesystem. In contexts requiring extremely robust
security, it may be best to consider another option and/or to consult a
professional.