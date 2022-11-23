package utils

import (
	"os"
	"sort"

	common "github.com/bozdoz/nevermind/nvm-common"
)

func GetInstalledVersions() (versions []common.Version, err error) {
	dir, err := common.GetNVMDir("node")

	if err != nil {
		return
	}

	files, err := os.ReadDir(dir)

	if err != nil {
		return
	}

	versions = make([]common.Version, 0, len(files))

	for _, f := range files {
		if f.IsDir() {
			ver, err := common.GetVersion(f.Name())

			if err == nil {
				versions = append(versions, ver)
			}
		}
	}

	// ">" gives us DESC sort
	sort.Slice(versions, func(i, j int) bool {
		vi := versions[i]
		vj := versions[j]

		if vi.Major() == vj.Major() {
			if vi.Minor() == vj.Minor() {
				return vi.Patch() > vj.Patch()
			}
			return vi.Minor() > vj.Minor()
		}
		return vi.Major() > vj.Major()
	})

	return
}
