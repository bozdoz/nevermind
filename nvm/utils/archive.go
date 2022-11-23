package utils

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

func UnArchiveBytes(b []byte, dir string) (err error) {
	reader := bytes.NewReader(b)

	return UnArchiveReader(reader, dir)
}

func UnArchiveReader(reader io.Reader, dir string) (err error) {
	gzReader, err := gzip.NewReader(reader)

	if err != nil {
		return
	}

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// prevent zip slip vulnerability
		dest := filepath.Join(dir, header.Name)
		if !strings.HasPrefix(dest, filepath.Clean(dir)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", header.Name)
		}

		if header.Typeflag == tar.TypeSymlink {
			log.Println("found symlink", header.Linkname, dest)
			syscall.Symlink(header.Linkname, dest)
			continue
		}

		info := header.FileInfo()

		if info.IsDir() {
			if err = os.MkdirAll(dest, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(dest, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())

		if err != nil {
			return err
		}

		defer file.Close()

		_, err = io.Copy(file, tarReader)

		if err != nil {
			return err
		}
	}

	return
}
