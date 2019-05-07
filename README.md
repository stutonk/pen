#### for unixes
`make && make install`

#### everybody else
`go build` and do what ye will with it

```
pen(1)                           USER COMMANDS                          pen(1)



NAME
       pen - encrypt and decrypt files with a password

SYNOPSIS
       pen [-hv] file [files...]

DESCRIPTION
       Pen  is  a lightweight alternatve to PGP for securely encrypting files;
       it uses the Argon2 algorithm for key-stretching  passwords  and  NaCl's
       secretbox for data encryption. If a filename ends in '.pen', it will be
       decrypted rather than encrypted. Upon successful encryption or  decryp‐
       tion, the original file or .pen file, respectively, is removed.

OPTIONS
       --help display this help and exit

       --version
              output version information and exit

SECURITY CONSIDERATIONS
       Pen, like almost any lock, will only stop honest people. Due to the na‐
       ture of underlying filesystems and their underlying  hardware,  no  at‐
       tempt is made to shred the original file before it is removed following
       encryption. A determined person, not even to  mention  law  enforcement
       agencies,  with  access  to  your  raw disk WILL be able to recover the
       original file. This proram is designed to be as simple and straightfor‐
       ward as possible so that it can be audited easily. However, the general
       assumption is also made that implementations of the underlying  crypto‐
       graphic  algorithms  are  both correct and secure. Therefore, it is ad‐
       vised that you use this program with caution, or not at all, for criti‐
       cal applications.

AUTHOR
       Joseph Eib (github.com/stutonk)



version 1.0.0                     7 May 2019                            pen(1)
```