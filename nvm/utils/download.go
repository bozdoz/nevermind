package utils

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

/*
actually downloads a url and passes the []bytes
to the channel;
if there is an error, we pass it to `err_ch`
used by [commands.Install]
*/
func Download(url string, ch chan []byte, err_ch chan error) {
	defer close(ch)
	s := time.Now()
	res, err := http.Get(url)

	if err != nil {
		err_ch <- err
		return
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		err_ch <- err
		return
	}

	res.Body.Close()

	log.Println("status", res.Status)
	log.Println("headers", res.Header)

	if res.StatusCode != 200 {
		err_ch <- fmt.Errorf("request of %s failed with status code: %d", url, res.StatusCode)
		return
	}

	log.Printf("downloaded %s (%s)", url, time.Since(s))

	ch <- body
}
