package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func SaveToFile(body []byte, filename string) error {
	file, err := os.Create(filename)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(body)

	if err != nil {
		removeErr := os.Remove(filename)

		if removeErr != nil {
			fmt.Println("filename failed to be removed:", filename)
		}

		return err
	}

	abs, _ := filepath.Abs(filename)

	log.Println("Downloaded", abs)

	return nil
}
