
package main

import (
    "fmt"
    "log"
	"net"
    "os"
    "flag"
    "golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
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
    Be careful use this for rm command or something like that
    Don't use this do harmful things
    `)
    os.Exit(1)
}

func main() {
	HOSTs := []string {
        "134.160.36.66",
        "134.160.36.67",
        "134.160.36.68",
        "134.160.36.69",
        "134.160.36.101",
        "134.160.36.102",
        "134.160.36.103",
        "134.160.36.104",
        "134.160.36.105",
        "134.160.36.106",
        "134.160.36.107",
        "134.160.36.108",
    }

	PASSs := []string {
        "asiainfo",
        "asiainfo",
        "asiainfo",
        "asiainfo",
        "bill01app",
        "bill01app",
        "bill01app",
        "bill01app",
        "bill01app",
        "bill01app",
        "bill01app",
        "bill01app",
    }
    flag.Parse()
    if os.Args == nil || shellCmd == "" {
        Usage()
    }

    response := make(chan string)
    //TODO: not dial on same machine

    USER := "bill01"
    for i, _:= range HOSTs {
        go dial(HOSTs[i],USER,PASSs[i],22,1<<15,shellCmd,response)
    }
    for j := 0; j<len(HOSTs);j++ {
        select {
        case res:= <-response:
            fmt.Println(res)
        }
    }
    close(response)
}

func dial(HOST string,USER string,PASS string,PORT int,SIZE int,shellCmd string,res chan string) {
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

    b, err := session.Output(shellCmd)
    if err != nil {
        log.Fatalf("failed to execute: %s", err)
    }
    log.Println("--------------------------------------")
    log.Println("Output: ",HOST)
    log.Println(string(b))
	res <- "done"+HOST
}
