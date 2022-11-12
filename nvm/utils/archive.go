package utils

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

// TODO: gzipping can be done with the tarring:
// https://cs.opensource.google/go/x/build/+/8c11e572:internal/untar/untar.go
// reverse the effects of tarring
func Untar(tarball, target string) error {
	reader, err := os.Open(tarball)
	if err != nil {
		return err
	}
	defer reader.Close()
	tarReader := tar.NewReader(reader)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		// TODO: zip slip vulnerability
		// destpath := filepath.Join(destination, filePath)
		// if !strings.HasPrefix(destpath, filepath.Clean(destination)+string(os.PathSeparator)) {
		// 	return fmt.Errorf("%s: illegal file path", filePath)
		// }

		path := filepath.Join(target, header.Name)

		info := header.FileInfo()

		if header.Typeflag == tar.TypeSymlink {
			log.Println("found symlink", header.Linkname, path)
			syscall.Symlink(header.Linkname, path)
			continue
		}
		if info.IsDir() {
			if err = os.MkdirAll(path, info.Mode()); err != nil {
				return err
			}
			continue
		}

		file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(file, tarReader)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnGzip(source, target string) (string, error) {
	reader, err := os.Open(source)
	if err != nil {
		fmt.Println("failed to open", source)
		return "", err
	}
	defer reader.Close()

	archive, err := gzip.NewReader(reader)
	if err != nil {
		return "", err
	}
	defer archive.Close()

	os.MkdirAll(target, 0755)

	target = filepath.Join(target, archive.Name)
	writer, err := os.Create(target)
	if err != nil {
		return "", err
	}
	defer writer.Close()

	_, err = io.Copy(writer, archive)

	return target, err
}
