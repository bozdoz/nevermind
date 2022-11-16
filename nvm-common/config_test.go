package common

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"
	"testing"
)

// mock a file as returned from [os.Open]
type mockReadWriteCloser struct {
	read  []byte
	write []byte
	err   error
}

func (mrc *mockReadWriteCloser) Read(p []byte) (n int, err error) {
	copy(p, mrc.read)
	return len(p), mrc.err
}

func (mrc *mockReadWriteCloser) Write(p []byte) (n int, err error) {
	mrc.write = append(mrc.write, p...)
	return len(p), mrc.err
}

func (mrc *mockReadWriteCloser) Close() error { return nil }

var errNoNvmPath = errors.New("no nvm path")
var errFailedOpen = errors.New("can't open file")
var errFailedWrite = errors.New("can't write file")

func TestGetConfig(t *testing.T) {
	file := &mockReadWriteCloser{
		read: []byte("{}"),
	}

	defaultFileOpener := func(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
		return file, nil
	}
	defaultNvmDir := func(path ...string) (string, error) {
		return filepath.Join("some", CONFIG_NAME), nil
	}

	t.Run("returns config", func(t *testing.T) {
		// monkey-patching
		getNvmDir = defaultNvmDir
		fileOpener = defaultFileOpener

		cfg, err := GetConfig()

		if err != nil {
			t.Errorf("expected no error, got: %q", err)
		}

		if cfg.Current != "" {
			t.Errorf("expected empty cfg, got: %q", cfg)
		}
	})

	t.Run("returns error if there's a problem with config path", func(t *testing.T) {
		// monkey-patching
		getNvmDir = func(path ...string) (string, error) {
			return "", errNoNvmPath
		}
		fileOpener = defaultFileOpener

		_, err := GetConfig()

		if err != errNoNvmPath {
			t.Errorf("expected %q, got: %q", errNoNvmPath, err)
		}
	})

	t.Run("returns error if there's a problem with open or create file", func(t *testing.T) {
		// monkey-patching
		getNvmDir = defaultNvmDir
		fileOpener = func(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
			return file, errFailedOpen
		}

		_, err := GetConfig()

		if err != errFailedOpen {
			t.Errorf("expected %q, got: %q", errFailedOpen, err)
		}
	})

	version := "18.0.0"
	t.Run("it parses json", func(t *testing.T) {
		hasVersion := &mockReadWriteCloser{
			read: []byte(fmt.Sprintf("{\"current\":\"%s\"}", version)),
		}

		fileOpener = func(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
			return hasVersion, nil
		}

		getNvmDir = defaultNvmDir

		cfg, err := GetConfig()

		if err != nil {
			t.Errorf("expected no error, got: %q", err)
		}

		if cfg.Current != Version(version) {
			t.Errorf("expected %q, got %q", version, cfg.Current)
		}
	})

	t.Run("it fails on unparsable json", func(t *testing.T) {
		hasBadVersion := &mockReadWriteCloser{
			read: []byte(fmt.Sprintf("{current:%s}", version)),
		}

		fileOpener = func(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
			return hasBadVersion, nil
		}

		getNvmDir = defaultNvmDir

		_, err := GetConfig()

		want := "invalid character"

		if !strings.HasPrefix(err.Error(), want) {
			t.Errorf("expected error: %q, got: %q", want, err)
		}
	})
}

func TestSetConfig(t *testing.T) {
	file := &mockReadWriteCloser{}
	cfg := config{
		Current: "16.1.0",
	}

	t.Run("sets version", func(t *testing.T) {
		fileOpener = func(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
			return file, nil
		}

		err := SetConfig(cfg)

		if err != nil {
			t.Errorf("no error expected; got: %q", err)
		}

		want := "{\"current\":\"16.1.0\"}\n"
		got := string(file.write)

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("fails if no file", func(t *testing.T) {
		fileOpener = func(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
			return nil, errFailedOpen
		}

		err := SetConfig(cfg)

		if !errors.Is(err, errFailedOpen) {
			t.Errorf("error expected: %q; got: %q", errFailedOpen, err)
		}
	})

	t.Run("fails if json encoder fails", func(t *testing.T) {
		expectedErr := errFailedWrite
		fileOpener = func(name string, flag int, perm fs.FileMode) (io.ReadWriteCloser, error) {
			return &mockReadWriteCloser{
				err: expectedErr,
			}, nil
		}

		err := SetConfig(cfg)

		if !errors.Is(err, expectedErr) {
			t.Errorf("error expected: %q; got: %q", expectedErr, err)
		}
	})
}
