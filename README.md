pen
===
Pen is a lightweight alternatve to PGP for securely encrypting files; it uses
the Argon2 algorithm for key-stretching passwords and NaCl's secretbox for
data encryption. If a filename ends in '.pen', it will be decrypted rather
than encrypted. Upon successful encryption or decryption, the original file
or .pen file, respectively, is removed.

### for unixes
`make && make install`

### everybody else
`go build` and do what ye will with it

### note
Pen, like almost any lock, will only stop honest people. Due to the nature
of underlying filesystems and their underlying hardware, no attempt is made
to shred the original file before it is removed following encryption. A
determined person, not even to mention law enforcement agencies, with access
to your raw disk WILL be able to recover the original file. This proram is
designed to be as simple and straightforward as possible so that it can be
audited easily. However, the general assumption is also made that
implementations of the underlying cryptographic algorithms are both correct
and secure. Therefore, it is advised that you use this program with caution,
or not at all, for critical applications.