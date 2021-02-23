package engine

import (
	"os"
	"strings"
)

var OSPS = string(os.PathSeparator)

func shouldExclude(dirPath, filePath string, excludes []string) bool {
	parentPath := strings.TrimPrefix(dirPath, "."+OSPS)
	for _, ePath := range excludes {
		var exd = filepath.Dir(parentPath + addDirSuffix(ePath))
		var fpd = filepath.Dir(filePath)
		if exd == fpd {
			return true
		}
	}
	return false
}

func pathCompatible(dirPath string) string {
	if strings.Index(dirPath, OSPS) != -1 {
		return dirPath
	}
	return strings.ReplaceAll(dir, "/", OSPS)
}

func addDirSuffix(dirPath string) string {
	if !strings.HasSuffix(dirPath, OSPS) {
		dirPath += OSPS
	}
	return dirPath
}

func path2ObjectKey(dir string) string {
	return strings.ReplaceAll(dir, OSPS, "/")
}

func objectKey2Path(dir string) string {
	return strings.ReplaceAll(dir, "/", OSPS)
}