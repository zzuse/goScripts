
package main

import (
    "golang.org/x/crypto/ssh"
    "fmt"
    "log"
    "flag"
//    "sync"
//    "strings"
//   "io"
    "os"
)

var (
    debugSwitch            bool
    iManUser, iManPassword string
    hostIp, port, userName, passWord, shellCmd string
)

func init() {
    flag.StringVar(&hostIp, "h", "", "The remote host ip.")
    flag.StringVar(&port, "port", "22", "The remote host port.")
    flag.StringVar(&userName, "u", "", "The user name login host.")
    flag.StringVar(&passWord, "p", "", "Log on to the host password.")
    flag.StringVar(&iManUser, "iManu", "", "The user name login iMan.")
    flag.StringVar(&iManPassword, "iManp", "", "Login iMan's password.")
    flag.StringVar(&shellCmd, "c", "", "To perform a shell command.")
    flag.BoolVar(&debugSwitch, "d", false, "The debug switch.")
}

func Usage() {
    fmt.Printf(`Usage of cssh:
    -c string
    To perform a shell command.
    -d bool The debug statu.
    -h string
    The remote host ip.
    -iManp string
    The host login iman user password.
    -iManu string
    The host login iman user name.
    -p string
    The host login user password.
    -port string
    The remote host port. (default "22")
    -u string
    The host login user name.
    [cssh -h 192.168.1.1 -port 22 -u root -p root -iManU admin -iManP admin -c "uptime;whoami" -d true]`)
    os.Exit(1)
}

func main() {
    flag.Parse()
    if os.Args == nil || hostIp == "" || userName == "" ||
    passWord == "" || shellCmd == "" {
        Usage()
    }

    config := &ssh.ClientConfig{
        User: userName,
        Auth: []ssh.AuthMethod{
            ssh.Password(passWord),
        },
    }
    c, err := ssh.Dial("tcp", hostIp+":"+port, config)
    if err != nil {
        log.Println("unable to dial remote side:", err)
    }
    defer c.Close()

    // Create a session
    session, err := c.NewSession()
    if err != nil {
        log.Fatalf("unable to create session: %s", err)
    }
    defer session.Close()

    b, err := session.Output(shellCmd)
    if err != nil {
        log.Fatalf("failed to execute: %s", err)
    }
    log.Println("Output: ", string(b))

    return
}
