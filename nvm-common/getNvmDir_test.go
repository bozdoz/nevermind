package common

import (
	"errors"
	"fmt"
	"path/filepath"
	"testing"
)

var mockHomeDir = filepath.Join(string(filepath.Separator), "home", "gopher")

var expectBase = filepath.Join(mockHomeDir, NVM_DIR)

var goodHomeGetter = func() (string, error) {
	return mockHomeDir, nil
}

var errNoHome = errors.New("no home")

var badHomeGetter = func() (string, error) {
	return "", errNoHome
}

func TestGetNodeBin(t *testing.T) {
	v := "18.0.0"
	bin := "yarn"
	goodBin := filepath.Join(expectBase, "node", v, "bin", bin)
	errNoFile := errors.New("no file")
	fileExists := func(path string) error {
		return nil
	}
	notFileExists := func(path string) error {
		return errNoFile
	}
	t.Run("can get package with version", func(t *testing.T) {
		file, err := getNodeBinWithGetter(goodHomeGetter, fileExists, Version(v), bin)

		if err != nil {
			t.Errorf("expected no error, but got: %s", err)
		}

		if file != goodBin {
			t.Errorf("got %q, want %q", file, goodBin)
		}
	})

	t.Run("handles no homedir", func(t *testing.T) {
		_, err := getNodeBinWithGetter(badHomeGetter, fileExists, Version(v), bin)

		if err == nil || !errors.Is(err, errNoHome) {
			t.Errorf("expected home error, but got: %s", err)
		}
	})

	t.Run("handles non-existent bin", func(t *testing.T) {
		_, err := getNodeBinWithGetter(goodHomeGetter, notFileExists, Version(v), bin)

		if err == nil || !errors.Is(err, errNoFile) {
			t.Errorf("expected no file error, but got: %s", err)
		}
	})
}

// this seems pretty contrived
func TestGetNVMDir(t *testing.T) {
	tests := [...][]string{
		{},
		{"node"},
		{"node", "1.8.0", "bin", "tsc"},
	}

	// adds path to expectBase
	prependBase := func(paths ...string) string {
		paths = append([]string{expectBase}, paths...)
		return filepath.Join(paths...)
	}

	for _, vals := range tests {
		t.Run(fmt.Sprintf("%q should return %q", vals, vals), func(t *testing.T) {
			v, err := getNVMDirWithGetter(goodHomeGetter, vals...)
			expected := prependBase(vals...)

			if err != nil {
				t.Errorf("expected an error-free experience, got %q", err)
			} else if v != expected {
				t.Errorf("got %q, want %q", v, expected)
			}
		})
	}

	t.Run("fails if home dir can't be found", func(t *testing.T) {
		_, err := getNVMDirWithGetter(badHomeGetter, tests[0]...)

		if err == nil {
			t.Errorf("expected an error-ridden experience, got nil somehow")
		}
	})
}
