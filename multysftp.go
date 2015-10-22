//User and Pass and IP need read from text
//go routine for different dirs
//different dir may have same IPs.
//Authorized by Zhang Zhen
package main

import (
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

func main() {
/*	os.walk("")
	for walk result {
	    if LOCALDIR == 01 then HOST USER PASS PORT REMOTEDIR setting
	    go dial(HOST USER PASS PORT LOCALDIR REMOTEDIR)
	}
*/
    LOCALDIRs := []string {
        "/bill02/predeal/gprs/out/gcdr3g_01/",
        "/bill02/predeal/gprs/out/gcdr3g_02/",
        "/bill02/predeal/gprs/out/gcdr3g_03/",
        "/bill02/predeal/gprs/out/gcdr3g_04/",
        "/bill02/predeal/gprs/out/gcdr3g_05/",
        "/bill02/predeal/gprs/out/gcdr3g_07/",
        "/bill02/predeal/gprs/out/gcdr3g_08/",
        "/bill02/predeal/gprs/out/gcdr3g_09/",
        "/bill02/predeal/gprs/out/gcdr3g_10/",
        "/bill02/predeal/gprs/out/gcdr3g_11/",
    }
    REMOTEDIRs := []string {
        "/bill02/predeal/gprs/out/gcdr3g_01/tmp/",
        "/bill02/predeal/gprs/out/gcdr3g_02/tmp/",
        "/bill02/predeal/gprs/out/gcdr3g_03/tmp/",
        "/bill02/predeal/gprs/out/gcdr3g_04/tmp/",
        "/bill02/predeal/gprs/out/gcdr3g_05/tmp/",
        "/bill02/predeal/gprs/out/gcdr3g_07/tmp/",
        "/bill02/predeal/gprs/out/gcdr3g_08/tmp/",
        "/bill02/predeal/gprs/out/gcdr3g_09/tmp/",
        "/bill02/predeal/gprs/out/gcdr3g_10/tmp/",
        "/bill02/predeal/gprs/out/gcdr3g_11/tmp/",
    }

    HOSTs := []string {
        "134.160.36.105",
        "134.160.36.105",
        "134.160.36.106",
        "134.160.36.106",
        "134.160.36.106",
        "134.160.36.107",
        "134.160.36.107",
        "134.160.36.107",
        "134.160.36.108",
        "134.160.36.108",
    }

    response := make(chan string)

    USER := "bill01"
    PASS := "bill01app"
    for i, LOCALDIR := range LOCALDIRs {
        go dial(HOSTs[i],USER,PASS,22,1<<15,LOCALDIR,REMOTEDIRs[i],response)
    }
    for j := 0; j<len(LOCALDIRs);j++ {
        select {
        case res:= <-response:
            fmt.Println(res)
        }
    }
    close(response)
}

func dial(HOST string,USER string,PASS string,PORT int,SIZE int,LOCALDIR string,REMOTEDIR string, res chan string) {
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

	c, err := sftp.NewClient(conn, sftp.MaxPacket(SIZE))
	if err != nil {
		log.Fatalf("unable to start sftp subsytem: %v", err)
	}
	defer c.Close()

    d, err := os.Open(LOCALDIR)
    if err!= nil {
		log.Fatalf("unable to open local dir : %v", err)
    }
    defer d.Close()
    for {
        fileList,err := d.Readdir(100)
        if err == io.EOF {
           break;
        }
        for _,readFile := range fileList {
            if readFile.IsDir() == true {
                continue;
            }
            log.Printf("writing name %s ", readFile.Name())
            f, err := os.Open(LOCALDIR+readFile.Name())
            if err != nil {
                log.Fatal(err)
            }
            defer f.Close()
            info, _ := f.Stat();

            outputFile:=REMOTEDIR+readFile.Name()
            w, err := c.OpenFile(outputFile, syscall.O_CREAT|syscall.O_TRUNC|syscall.O_RDWR)
            if err != nil {
                log.Fatal(err)
            }
            defer w.Close()
            //log.Printf("wrote aa " )

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
            os.Remove(LOCALDIR+readFile.Name())
            log.Printf("local file removed %s", readFile.Name())
        }
    }
    res <- "done"+LOCALDIR
}

