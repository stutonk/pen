package main

import (
	"fmt"

	"github.com/spf13/pflag"
	"golang.org/x/crypto/ssh/terminal"
)

var decrypt bool

func init() {
	pflag.BoolVarP(&decrypt, "decrypt", "d", false, "switch to decrypt mode")
	pflag.Parse()
}

func main() {}

func promptForPass() []byte {
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
