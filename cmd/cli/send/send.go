package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

var (
	key  string
	cli1 string
	cli2 string
	name string
)

func main() {
	if len(os.Args) <= 1 {
		loadHelp()
		return
	}
	argsWithoutProg := os.Args[1:]
	if argsWithoutProg[0] == "help" || argsWithoutProg[0] == "-h" {
		loadHelp()
		return
	}
	hasKey := false
	size := len(argsWithoutProg)
	for i, v := range argsWithoutProg {
		if v == "-k" {
			hasKey = true
			if i < size-1 {
				key = argsWithoutProg[i+1]
				argsWithoutProg = append(argsWithoutProg[:i], argsWithoutProg[i+2:]...)
			} else {
				fmt.Fprintln(os.Stderr, "no private key found")
				return
			}
		}
	}

	if runtime.GOOS == "windows" {
		cli1 = "cli.exe"
		cli2 = "chain33-cli.exe"
	} else {
		cli1 = "cli"
		cli2 = "chain33-cli"
	}

	_, err := os.Stat(cli1)
	if err == nil {
		name = "cli"
	}
	if os.IsNotExist(err) {
		_, err = os.Stat(cli2)
		if err == nil {
			name = "chain33-cli"
		}
		if os.IsNotExist(err) {
			fmt.Println("no compiled cli file found")
			return
		}
	}

	cmdCreate := exec.Command(name, argsWithoutProg...)
	var outCreate bytes.Buffer
	var errCreate bytes.Buffer
	cmdCreate.Stdout = &outCreate
	cmdCreate.Stderr = &errCreate
	err = cmdCreate.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if errCreate.String() != "" {
		fmt.Println(errCreate.String())
		return
	}
	//fmt.Println("unsignedTx", outCreate.String(), errCreate.String())

	if !hasKey || key == "" {
		fmt.Fprintln(os.Stderr, "no private key found")
		return
	}
	bufCreate := outCreate.Bytes()
	cmdSign := exec.Command(name, "wallet", "sign", "-d", string(bufCreate[:len(bufCreate)-1]), "-k", key)
	var outSign bytes.Buffer
	var errSign bytes.Buffer
	cmdSign.Stdout = &outSign
	cmdSign.Stderr = &errSign
	err = cmdSign.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if errSign.String() != "" {
		fmt.Println(errSign.String())
		return
	}
	//fmt.Println("signedTx", outSign.String(), errSign.String())

	bufSign := outSign.Bytes()
	cmdSend := exec.Command(name, "tx", "send", "-d", string(bufSign[:len(bufSign)-1]))
	var outSend bytes.Buffer
	var errSend bytes.Buffer
	cmdSend.Stdout = &outSend
	cmdSend.Stderr = &errSend
	err = cmdSend.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	if errSend.String() != "" {
		fmt.Println(errSend.String())
		return
	}
	bufSend := outSend.Bytes()
	fmt.Println(string(bufSend[:len(bufSend)-1]))
}

func loadHelp() {
	fmt.Println("Use similarly as bty/token/trade raw transaction creation, in addition to the parameter of private key input following \"-k\".")
}