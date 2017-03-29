package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	_ "time"
)

const (
	bufferSize = 1280 * 1024 //写图片文件的缓冲区大小
)

var (
	numPoller = flag.Int("p", 1, "page loader num")
	//numDownloader = flag.Int("d", 0, "image downloader num")
	savePath = flag.String("s", "./downloads/", "save path")
	//imgExp        = regexp.MustCompile(`<a\s+class="img"\s+href="[a-zA-Z0-9_\-/:\.%?=]+">[\r\n\t\s]*<img\s+src="([^"'<>]*)"\s*/?>`)
	img2Exp = regexp.MustCompile(`<a href="(.*)" class="download-link">`)
)

type image struct {
	url      string
	filename string
}

type sexyContext struct {
	pollerDone   chan struct{}
	images       map[string]int
	imagesLock   *sync.Mutex
	imageChan    chan *image
	pageIndex    int32
	rootURL      string
	done         bool
	imageCounter int32
	okCounter    int32
}

func main() {
	flag.Parse()
	ctx := &sexyContext{
		pollerDone: make(chan struct{}),
		images:     make(map[string]int),
		imagesLock: &sync.Mutex{},
		imageChan:  make(chan *image, 100),
		pageIndex:  1,
		//rootURL:    "http://cuoss.asiainfo.com/cgi-bin/cvsweb/cvsweb.cgi/products/unibss/binary/license/Attic/hunan.lic",
		rootURL: "http://cuoss.asiainfo.com/cgi-bin/cvsweb/cvsweb.cgi/products/unibss/binary/license/Attic/yunnan.lic",
	}
	os.MkdirAll(*savePath, 0777)
	ctx.start()

}

func (ctx *sexyContext) start() {
	fmt.Printf("Poller%d\n", *numPoller)
	for i := 0; i < *numPoller; i++ {
		go ctx.downloadPage()
	}
	//fmt.Printf("download%d\n", *numDownloader)
	waits := sync.WaitGroup{}

	<-ctx.pollerDone
	ctx.done = true
	//close(ctx.pollerDone)
	waits.Wait()
	fmt.Printf("fetch done get img %d lic ok %d\n", ctx.imageCounter, ctx.okCounter)
}

func (ctx *sexyContext) downloadPage() {
	isDone := false
	for !isDone {
		select {
		case <-ctx.pollerDone:
			isDone = true
		default:
			url := fmt.Sprintf("%s", ctx.rootURL)
			fmt.Printf("download page %s\n", url)
			resp, err := http.Get(url)
			if err != nil {
				fmt.Printf("failed to load url %s with error %v", url, err)
			} else {
				defer resp.Body.Close()
				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					fmt.Printf("failed to load url %s with error %v", url, err)
				} else {
					ctx.parsePage(body)
				}
			}
		}
	}
}

func (ctx *sexyContext) parsePage(body []byte) {
	//fmt.Printf("%s\n", string(body))
	body2 := string(body)
	//body is like <a href="/cgi-bin/cvsweb/cvsweb.cgi/~checkout~/products/unibss/binary/license/Attic/yunnan.lic?rev=1.1.2.30;content-type=application%2Foctet-stream" class="download-link">download</a>
	idx := img2Exp.FindAllStringSubmatch(body2, -1)
	if idx == nil {
		ctx.pollerDone <- struct{}{}
	} else {
		fmt.Printf("%d\n", len(idx))
		for _, n := range idx {
			url := fmt.Sprintf("http://cuoss.asiainfo.com%s", n[1])
			str := strings.Split(url, "/")
			length := len(str)
			imgeUrl := url
			//get filename by "?"
			tmpfilename := strings.Split(str[length-1], "?")
			filename := tmpfilename[0]
			image := &image{url: imgeUrl, filename: filename}
			//atomic.AddInt32(&ctx.imageCounter, 1)
			//ctx.imageChan <- image
			fmt.Printf("start download %s\n", image.url)
			atomic.AddInt32(&ctx.okCounter, 1)
			if ctx.okCounter > 1 {
				ctx.pollerDone <- struct{}{}
				fmt.Printf("counter greater than 1 abort download %s\n", image.url)
				break
			}
			resp, err := http.Get(image.url)
			if err != nil {
				fmt.Printf("failed to load url %s with error %v\n", image.url, err)
			} else {
				defer resp.Body.Close()
				saveFile := *savePath + image.filename //path.Base(imgUrl)

				img, err := os.Create(saveFile)
				if err != nil {
					fmt.Print(err)

				} else {
					defer img.Close()

					imgWriter := bufio.NewWriterSize(img, bufferSize)

					_, err = io.Copy(imgWriter, resp.Body)
					if err != nil {
						fmt.Print(err)

					}
					imgWriter.Flush()
					fmt.Printf("finish download %s\n", image.url)
				}
			}
		}
	}
}
