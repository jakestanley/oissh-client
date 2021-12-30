package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"strings"

	ui "github.com/gizak/termui/v3"
	"golang.org/x/crypto/ssh"
)

var (
	hasSession bool = false
	client     *ssh.Client
	session    *ssh.Session
	stdin      io.WriteCloser
)

func main() {

	// for vscode debugging purposes
	// connect("localhost:2223")

	initUi()
	defer ui.Close()
	go renderUi()

	for inputUi() {
	}

	if hasSession {
		defer session.Close()
		defer client.Close()
	}

	// TODO fix. segfaults
	// err := session.Close()
	// if err != nil {
	// 	log.Println("Could not close session. Probably not open or already closed")
	// } else {
	// 	log.Println("SSH session ended gracefully")
	// }
}

func getPublicKey() ssh.AuthMethod {
	// publicKeyPath := "/Users/jake/.ssh/oissh.pub"
	privateKeyPath := "/Users/jake/.ssh/oissh"

	buffer, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal(err)
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		log.Fatal(err)
	}

	return ssh.PublicKeys(key)
}

func cmd(cmd string) {

	fmt.Printf("Sending: '%s'", cmd)
	fmt.Fprintf(stdin, "%s\n", cmd)

	// if err := session.Start(cmd); err != nil {
	// 	log.Println("session.Run error")
	// 	log.Fatal(err)
	// }
}

func connect(cxn string) {

	fmt.Printf("Connecting to %s\n", cxn)
	ppk := getPublicKey()
	user := "jake"

	verificationCallback := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		// TODO implement
		return ssh.InsecureIgnoreHostKey()(hostname, remote, key)
		// this will be what confirms the server identity,
		// 	like that MitM message you get when a server
		// 	you're connecting to via CLI has its fingerprint changed
	}

	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ppk,
		},
		HostKeyCallback: verificationCallback,
	}

	var err error
	client, err = ssh.Dial("tcp", cxn, config)
	if err != nil {
		log.Fatal(err)
	}

	session, err = client.NewSession()
	if err != nil {
		log.Fatal(err)
	}
	hasSession = true

	// modes := ssh.TerminalModes{
	// 	ssh.ECHO:  0,
	// 	ssh.IGNCR: 1,
	// }

	// if err = session.RequestPty("vt100", 80, 40, modes); err != nil {
	// 	log.Println("session.RequestPty error")
	// 	log.Fatal(err)
	// }

	stdin, err = session.StdinPipe()
	if err != nil {
		log.Println("session.StdinPipe error")
		log.Fatal(err)
	}

	if err = session.Shell(); err != nil {
		log.Println("session.Shell error")
		log.Fatal(err)
	}

	log.Println("Connected!")
}

func processInput(input string) {
	commands := strings.Fields(input)

	if len(commands) == 0 {
		return
	}

	switch strings.ToLower(commands[0]) {
	case "c", "connect":
		if commands[1] == "default" {
			go connect("localhost:2223")
		} else {
			go connect(commands[1])
		}

	// case "hello":
	//
	default:
		cmd(input)
		return
	}
}
