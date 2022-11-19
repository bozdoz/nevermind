package utils

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	common "github.com/bozdoz/nevermind/nvm-common"
)

// version is non-specific and we need to query the releases
// to find the latest version which fits the range
func GetLTS() (out common.Version, err error) {
	s := time.Now()

	defer log.Println("getLTS time", time.Since(s))

	url, err := url.JoinPath(BASE_URL, INDEX)

	if err != nil {
		return
	}

	res, err := http.Get(url)

	if err != nil {
		return
	}

	scanner := bufio.NewScanner(res.Body)

	var lts_index int

	if scanner.Scan() {
		text := scanner.Text()
		// gets header row
		fields := strings.Fields(text)

		for i, header := range fields {
			if header == "lts" {
				lts_index = i
			}
		}
	}

	for scanner.Scan() {
		text := scanner.Text()
		fields := strings.Fields(text)

		// check for first version that doesn't have "-" under lts
		if fields[lts_index] != "-" {
			out, err = common.GetVersion(fields[0])
			log.Println("lts", out, fields[1], fields[lts_index])

			return
		}
	}

	return
}
