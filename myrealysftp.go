// streaming-write-benchmark benchmarks the peformance of writing
// from /dev/zero on the client to /dev/null on the server via io.Copy.
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"

	"github.com/pkg/sftp"
)

var (
	//USER = flag.String("user", os.Getenv("USER"), "ssh username")
	USER = flag.String("user", "tiansl", "ssh username")
	HOST = flag.String("host", "192.130.12.38", "ssh server hostname")
	PORT = flag.Int("port", 22, "ssh server port")
	PASS = flag.String("p", os.Getenv("SOCKSIE_SSH_PASSWORD"), "ssh password")
	SIZE = flag.Int("s", 1<<15, "set max packet size")
	FILE = flag.String("f", "a.out", "upload filename")
	DIR  = flag.String("d", "./", "upload dirname")
)

func init() {
	flag.Parse()
}

func main() {
	var auths []ssh.AuthMethod
	if aconn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		auths = append(auths, ssh.PublicKeysCallback(agent.NewClient(aconn).Signers))

	}
	if *PASS != "" {
		auths = append(auths, ssh.Password(*PASS))
	}

	config := ssh.ClientConfig{
		User: *USER,
		Auth: auths,
	}
	addr := fmt.Sprintf("%s:%d", *HOST, *PORT)
	conn, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		log.Fatalf("unable to connect to [%s]: %v", addr, err)
	}
	defer conn.Close()

	c, err := sftp.NewClient(conn, sftp.MaxPacket(*SIZE))
	if err != nil {
		log.Fatalf("unable to start sftp subsytem: %v", err)
	}
	defer c.Close()

    //TODO: a bunch of local dirs
    //maybe need regexp to match Files
    d, err := os.Open(*DIR)
    if err!= nil {
		log.Fatalf("unable to open local dir : %v", err)
    }
    defer d.Close()
    fileList,_ := d.Readdir(100)
    for _,readFile := range fileList {
        if readFile.IsDir() == true {
            continue;
        }
        if readFile.Name() != *FILE {
            continue;
        }
        log.Printf("writing name %s ", readFile.Name())
        f, err := os.Open(readFile.Name())
        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()
        info, _ := f.Stat();

        //TODO: replace hard code dir to a routed configured dir 
        outputFile:="/unibss/tstusers/tiansl/zhangzhen/"+readFile.Name()
        w, err := c.OpenFile(outputFile, syscall.O_CREAT|syscall.O_TRUNC|syscall.O_RDWR)
        if err != nil {
            log.Fatal(err)
        }
        defer w.Close()

        const size int64 = 1e9

        log.Printf("writing %v bytes", info.Size())
        t1 := time.Now()
        n, err := io.Copy(w, io.LimitReader(f, info.Size()))
        if err != nil {
            log.Fatal(err)
        }
        if n != info.Size() {
            log.Fatalf("copy: expected %v bytes, got %d", info.Size(), n)
        }
        log.Printf("wrote %v bytes in %s", info.Size(), time.Since(t1))
    }
}
