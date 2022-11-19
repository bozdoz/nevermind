package utils

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	common "github.com/bozdoz/nevermind/nvm-common"
)

const INDEX = "index.tab"

// version is non-specific and we need to query the releases
// to find the latest version which fits the range
func GetLatestFromVersion(version common.Version) (out common.Version, err error) {
	s := time.Now()

	defer log.Println("getLatestVersion time", time.Since(s))

	url, err := url.JoinPath(BASE_URL, INDEX)

	if err != nil {
		return
	}

	res, err := http.Get(url)

	if err != nil {
		return
	}

	scanner := bufio.NewScanner(res.Body)

	for scanner.Scan() {
		text := scanner.Text()
		// check for "v{version}."
		if strings.HasPrefix(text, fmt.Sprintf("v%s.", version)) {
			text := strings.Split(text, "\t")[0]
			out, err = common.GetVersion(text)

			return
		}
	}

	return
}
