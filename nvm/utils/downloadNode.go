package utils

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// CI usually sets CI=true? disable progress bars?
var CI = os.Getenv("CI")

func DownloadNode(url string) (b []byte, err error) {
	// check valid url
	url, err = ValidateUrl(url)

	if err != nil {
		return
	}

	// time request
	s := time.Now()

	res, err := http.Get(url)

	log.Printf("%s request time (%s)", url, time.Since(s))

	if err != nil {
		return
	}

	defer res.Body.Close()

	log.Println("status", res.Status)
	log.Println("headers", res.Header)

	size := int(res.ContentLength)

	// time download
	s = time.Now()

	if size == -1 || CI != "" {
		// no progress bar possible without download size
		return io.ReadAll(res.Body)
	}

	// progressbar download
	b = make([]byte, 0, size)

	go ProgressBar(&b, size)

	// rewritten from io.ReadAll
	for {
		// read into b after length of b until end of b
		n, err := res.Body.Read(b[len(b):size])
		// grow length of b by n
		b = b[:len(b)+n]

		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
	}

	if err == nil {
		log.Println("node download time", time.Since(s))
	}

	return
}
