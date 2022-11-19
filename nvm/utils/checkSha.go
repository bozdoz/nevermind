package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
)

// error when verified sha did not match calculated file
var ErrInvalidSha = errors.New("invalid sha")

// gets file sha and checks against a verified sha
// returns [ErrInvalidSha]
func CheckSha(node_body []byte, sha string) (err error) {
	h := sha256.New()
	h.Write(node_body)
	file_sha := fmt.Sprintf("%x", h.Sum(nil))

	log.Println("verified sha", sha)
	log.Println("file sha", file_sha)

	if sha != file_sha {
		err = ErrInvalidSha
	}

	return
}
