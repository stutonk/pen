package main

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"

	flag "github.com/spf13/pflag"
	"github.com/stutonk/boxutil"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	cryptExt        = ".pen"
	errFmt          = "%v: (%v) fatal; %v\n"
	magic    uint32 = 0xc0ffee11
	saltLen  uint32 = 128
	usageFmt        = "usage: %v [-h, -v] file [files...]\nOptions are:\n"
	verFmt          = "%v version %v\n"
	version         = "1.1.0"
)

var (
	appName     = os.Args[0]
	helpFlag    bool
	order       = binary.BigEndian
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
file is encrypted.

The header format for .pen files is:
	magic   uint32
	saltLen uint32
	salt    [saltLen]byte
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

	pass := promptForPass("Enter password: ")

	defer func() {
		if r := recover(); r != nil {
			switch err := r.(type) {
			case fileErr:
				fmt.Printf(errFmt, appName, err.File, err.Err)
			default:
				fmt.Printf("%v: fatal; %v\n", appName, err)
			}
		}
	}()

	var confirmed bool
	for _, fileName := range flag.Args() {
		in, err := os.Open(fileName)
		if err != nil {
			panic(fileErr{fileName, err})
		}
		defer in.Close()

		var (
			op          func(io.Reader, io.Writer, *[32]byte) error
			outName     string
			salt        []byte
			writeHeader bool
		)

		if filepath.Ext(fileName) == cryptExt {
			var (
				hMagic   uint32
				hSaltLen uint32
			)
			if err := binary.Read(in, order, &hMagic); err != nil {
				panic(fileErr{fileName, err})
			}
			if hMagic != magic {
				panic(fileErr{
					fileName,
					fmt.Errorf("wrong magic number in header"),
				})
			}
			if err := binary.Read(in, order, &hSaltLen); err != nil {
				panic(fileErr{fileName, err})
			}
			salt = make([]byte, hSaltLen)
			if _, err := in.Read(salt); err != nil {
				panic(fileErr{fileName, err})
			}
			op = boxutil.OpenStream
			outName = fileName[:len(fileName)-4]
		} else {
			if !confirmed {
				repeat := promptForPass("Enter password (repeat): ")
				if !reflect.DeepEqual(repeat, pass) {
					panic(fmt.Errorf("passwords don't match"))
				}
				confirmed = true
			}
			salt = make([]byte, saltLen)
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
			binary.Write(out, order, magic)
			binary.Write(out, order, saltLen)
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

func promptForPass(msg string) []byte {
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
	fmt.Print(msg)
	pass, err := terminal.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	return pass
}
