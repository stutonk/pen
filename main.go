package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	flag "github.com/spf13/pflag"
	"github.com/stutonk/pen/pkg/boxutil"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	cryptExt        = ".pen"
	errFmt          = "%v: (%v) fatal; %v\n"
	magic    uint32 = 0x1c0ffee9
	saltLen         = 128
	usageFmt        = "usage: %v [-hv] file [files...]\n\n"
	verFmt          = "%v version %v\n"
	version         = "1.0.0"
)

var (
	appName     = os.Args[0]
	headerLen   = 4 + len(version) + saltLen
	helpFlag    bool
	verFlag     bool
	verValue, _ = strconv.ParseFloat(version[:3], 64)
)

type fileErr struct {
	File string
	Err  error
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usageFmt, appName)
		flag.PrintDefaults()
		fmt.Println()
	}

	flag.BoolVarP(
		&helpFlag,
		"help",
		"h",
		false,
		"display this help and exit",
	)
	flag.BoolVar(
		&verFlag,
		"version",
		false,
		"output version information and exit",
	)
	flag.Parse()
}

/*
Decryption is attempted when a file has a .pen extension, otherwise the
file is encrypted. Version compatibility is tested (as a float64) against
the MAJOR and MINOR elements of a semver version number.

The header format for .pen files is:
	magic   uint32
	version [5]byte
	salt    [128]byte
*/
func main() {
	switch {
	case helpFlag:
		flag.Usage()
		return
	case verFlag:
		fmt.Printf(verFmt, appName, version)
		return
	case len(flag.Args()) == 0:
		flag.Usage()
		return
	}

	pass := promptForPass()

	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case fileErr:
				fmt.Printf(errFmt, appName, err.File, err.Err)
			default:
				fmt.Printf("%v: fatal; %v", appName, err)
			}
		}
	}()

	for _, fileName := range flag.Args() {
		in, err := os.Open(fileName)
		if err != nil {
			panic(fileErr{fileName, err})
		}
		defer in.Close()

		salt := make([]byte, saltLen)
		var (
			op          func(io.Reader, io.Writer, *[32]byte) error
			outName     string
			writeHeader bool
		)

		if filepath.Ext(fileName) == cryptExt {
			header := make([]byte, headerLen)
			n, err := in.Read(header)
			switch {
			case err != nil && err != io.EOF:
				panic(fileErr{fileName, err})
			case n != headerLen:
				panic(fileErr{fileName, fmt.Errorf("invalid .pen file length")})
			case binary.BigEndian.Uint32(header[:4]) != magic:
				panic(fileErr{fileName, fmt.Errorf("wrong magic number in header")})
			}
			fileVer, err := strconv.ParseFloat(string(header[4:7]), 64)
			if err != nil {
				panic(fileErr{fileName, fmt.Errorf("corrupt version in header")})
			} else if verValue < fileVer {
				panic(fileErr{
					fileName,
					fmt.Errorf("need program version >= %.1f", fileVer),
				})
			}
			copy(salt, header[4+len(version):])
			op = boxutil.OpenStream
			outName = fileName[:len(fileName)-4]
		} else {
			if _, err := rand.Read(salt); err != nil {
				panic(fileErr{fileName, err})
			}
			op = boxutil.SealStream
			outName = fileName + cryptExt
			writeHeader = true
		}

		key := boxutil.Passkey(pass, salt)
		out, err := os.Create(outName)
		if err != nil {
			panic(fileErr{fileName, err})
		}
		defer out.Close()

		if writeHeader {
			binary.Write(out, binary.BigEndian, magic)
			out.Write([]byte(version))
			out.Write(salt)
		}

		if err = op(in, out, key); err != nil {
			defer os.Remove(outName)
			panic(fileErr{fileName, err})
		}

		if err = os.Remove(fileName); err != nil {
			defer os.Remove(outName)
			panic(fileErr{fileName, err})
		}
	}
}

func promptForPass() []byte {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("%v: fatal; %v", appName, r)
		}
	}()

	oldState, err := terminal.MakeRaw(0)
	if err != nil {
		panic(err)
	}
	defer func() {
		terminal.Restore(0, oldState)
		fmt.Println()
	}()
	fmt.Print("enter password: ")
	pass, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	return pass
}
