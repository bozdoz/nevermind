package utils

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"
)

func CheckSha(fileName string, node_body, sha_body []byte) error {
	h := sha256.New()
	h.Write(node_body)
	file_sha := fmt.Sprintf("%x", h.Sum(nil))

	log.Println("looking for", file_sha)

	shas := strings.Fields(string(sha_body))

	for i, line := range shas {
		if line == fileName {
			verified_sha := shas[i-1]
			log.Println("found", verified_sha)

			if verified_sha != file_sha {
				return fmt.Errorf("no SHA match for %s", fileName)
			} else {
				log.Println("SHA's match!")

				return nil
			}
		}
	}

	return fmt.Errorf("could not find filename (%s) in %s", fileName, SHASUMS)
}
