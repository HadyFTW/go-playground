package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

func main() {
	start := time.Now()
	ch := make(chan string)

	for _, uri := range os.Args[1:] {
		if !(strings.HasPrefix(uri, "http://") || strings.HasPrefix(uri, "https://")) {
			uri = fmt.Sprintf("%s%s", "http://", uri)
		}
		go fetch(uri, ch)
	}

	for range os.Args[1:] {
		fmt.Println(<-ch)
	}

	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}

func fetch(uri string, ch chan<- string) {
	start := time.Now()
	resp, err := http.Get(uri)
	defer resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}

	f, err := os.Create(url.QueryEscape(uri))
	defer f.Close()
	if err != nil {
		ch <- fmt.Sprint(err)
	}

	nbytes, err := io.Copy(f, resp.Body)
	if err != nil {
		ch <- fmt.Sprintf("while reading %s: %v", uri, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs\t%7d\t%s", secs, nbytes, uri)
}
