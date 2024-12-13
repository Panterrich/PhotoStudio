package rawmove

import (
	"path"
	"regexp"
	"strings"
)

const (
	RawDir = "RAW"
)

func isRawFile(fileName string) bool {
	rawFileExtension := map[string]struct{}{
		".cr2": {}, // Canon
		".nef": {}, // Nikon
		".arw": {}, // Sony
	}

	fileName = strings.ToLower(fileName)
	_, ok := rawFileExtension[path.Ext(fileName)]

	return ok
}

func modifyFileName(fileName string, camera string) string {
	r := regexp.MustCompile(`.*(?P<index>\d\d\d\d)\s*(\((?P<version>\d+)\))?\.(?P<ext>\w+)`)
	template := "${index}_${version}.${ext}"

	submatch := r.FindStringSubmatchIndex(fileName)
	if submatch == nil {
		return ""
	}

	result := []byte{}
	result = r.ExpandString(result, template, fileName, submatch)

	index := strings.ReplaceAll(string(result), "_.", ".")

	return strings.Join([]string{camera, index}, "_")
}
