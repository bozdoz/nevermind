package utils

import (
	"archive/tar"
	"archive/zip"
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
			_ = syscall.Symlink(header.Linkname, dest)
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

func UnZipBytes(b []byte, dir string) (err error) {
	zipReader, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))

	if err != nil {
		return err
	}

	for _, f := range zipReader.File {
		fpath := filepath.Join(dir, f.Name)

		// Prevent Zip Slip vulnerability
		if !strings.HasPrefix(fpath, filepath.Clean(dir)+string(os.PathSeparator)) {
			return &os.PathError{Op: "extract", Path: fpath, Err: os.ErrInvalid}
		}

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, os.ModePerm); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
