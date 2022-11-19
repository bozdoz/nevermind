package utils

import (
	"bufio"
	"errors"
	"log"
	"net/http"
	"strings"
)

// error when filename/sha not found
var ErrShaNotFound = errors.New("sha not found")

// downloads shasums and gets verified sha for a filename
func FetchSha(url, filename string) (sha string, err error) {
	url, err = ValidateUrl(url)

	if err != nil {
		return
	}

	log.Println("fetching sha from", url)

	res, err := http.Get(url)

	if err != nil {
		return
	}

	defer res.Body.Close()

	// read res.Body line-by-line
	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		// when we find the filename, we can end reading
		if fields[1] == filename {
			log.Println("got sha", fields[0])
			return fields[0], nil
		}
	}

	return "", ErrShaNotFound
}
