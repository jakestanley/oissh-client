package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	ui "github.com/gizak/termui/v3"
	"golang.org/x/crypto/ssh"
)

func main() {

	initUi()
	defer ui.Close()
	go renderUi()

	for inputUi() {
	}
}

func getPublicKey() string {
	publicKeyPath := "/Users/jake/.ssh/oissh.pub"
	privateKeyPath := "/Users/jake/.ssh/oissh"

	buffer, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	key, err = ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatal(err)
	}

	return ssh.PublicKeys(key), nil
}

func connect(cxn string) {

	ppk := getPublicKey()

	fmt.Printf("Connecting to %s\n", cxn)
}

func processInput(input string) {
	commands := strings.Fields(input)

	if len(commands) == 0 {
		return
	}

	switch strings.ToLower(commands[0]) {
	case "c", "connect":
		connect(commands[1])
	default:
		return
	}
}
