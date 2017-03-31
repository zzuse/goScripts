package main

import (
	//"bufio"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"log"
	"net"
	"os"
)

var (
	shellCmd string
)

func init() {
	flag.StringVar(&shellCmd, "c", "", "To perform a shell command.")
}

func Usage() {
	fmt.Printf(`Usage of multissh:
    -c string
    To perform a shell command on all the blade
    Be careful use this for rm command or something like that.
    `)
	fmt.Println("Don't use this do harmful things")
	os.Exit(1)
}

func main() {
	//here need do something nusty configure like json.
	HOSTs := []string{
		"135.64.20.143",
		"135.64.20.144",
		"135.64.20.145",
	}

	PASSs := []string{
		"asiainfo",
		"asiainfo",
		"asiainfo",
	}
	flag.Parse()
	if os.Args == nil || shellCmd == "" {
		Usage()
	}

	response := make(chan string)
	//TODO: not dial on same machine

	USER := "bill01"
	for i, _ := range HOSTs {
		go dial(HOSTs[i], USER, PASSs[i], 22, 1<<15, shellCmd, response)
	}
	for j := 0; j < len(HOSTs); j++ {
		select {
		case res := <-response:
			fmt.Println(res)
		}
	}
	close(response)
}

func dial(HOST string, USER string, PASS string, PORT int, SIZE int, shellCmd string, res chan string) {
	var auths []ssh.AuthMethod
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))
	}
	if PASS != "" {
		auths = append(auths, ssh.Password(PASS))
	}
	config := ssh.ClientConfig{
		User: USER,
		Auth: auths,
	}
	addr := fmt.Sprintf("%s:%d", HOST, PORT)
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		log.Fatalf("unable to connect to [%s]: %v", addr, err)
	}
	defer conn.Close()

	// Create a session
	session, err := conn.NewSession()
	if err != nil {
		log.Fatalf("unable to create session: %s", err)
	}
	defer session.Close()
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	running := true
	for running {
		log.Println("--------------------------------------")
		log.Println("HOST: ", HOST)
		log.Println(shellCmd)
		//session.Shell()
		//session.Wait()
		//b, err := session.Output(shellCmd)

		err := session.Run(shellCmd)
		if err != nil {
			log.Fatalf("failed to execute: %s", err)
		}
		shellCmd = string("exit")
		log.Println("HOST: ", shellCmd, HOST)
		log.Println("--------------------------------------")
		if shellCmd == "exit" {
			running = false
		}
	}
	res <- "done" + HOST
}
