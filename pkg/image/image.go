package image

import (
	"path"
	"strings"
)

func IsRaw(filename string) bool {
	rawFileExtension := map[string]struct{}{
		".cr2": {}, // Canon
		".cr3": {}, // Canon
		".nef": {}, // Nikon
		".arw": {}, // Sony
	}

	filename = strings.ToLower(filename)
	_, ok := rawFileExtension[path.Ext(filename)]

	return ok
}

func IsJpeg(filename string) bool {
	rawFileExtension := map[string]struct{}{
		".jpg":  {},
		".jpeg": {},
	}

	filename = strings.ToLower(filename)
	_, ok := rawFileExtension[path.Ext(filename)]

	return ok
}
