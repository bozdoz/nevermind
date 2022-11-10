package utils

import (
	"fmt"
	"io"
	"os"
)

func Rename(source, target string) (err error) {
	err = os.Rename(source, target)

	if err != nil {
		// try moving the file
		// could be an issue with running this
		// inside of a docker container
		err = moveFile(source, target)
	}

	return
}

func moveFile(source, target string) error {
	inputFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("couldn't open source file: %w", err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("couldn't open dest file: %w", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)

	if err != nil {
		return fmt.Errorf("writing to output file failed: %w", err)
	}

	err = os.Remove(source)
	if err != nil {
		return fmt.Errorf("failed removing original file: %w", err)
	}

	return nil
}