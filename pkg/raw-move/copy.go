package rawmove

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	log "github.com/rs/zerolog/log"
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

func copyFile(srcPath, dstDir string) error {
	var srcFileName, camera, dstPath string

	srcFileName = filepath.Base(srcPath)
	camera = filepath.Base(dstDir)

	srcFileName = modifyFileName(srcFileName, camera)
	if srcFileName == "" {
		return fmt.Errorf("empty filename: %s", srcFileName)
	}

	if isRawFile(srcFileName) {
		dstPath = filepath.Join(dstDir, RawDir, srcFileName)
	} else {
		dstPath = filepath.Join(dstDir, srcFileName)
	}

	if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
		log.Warn().Msgf("файл %s уже существует", dstPath)
		return nil
	}

	input, err := os.ReadFile(srcPath)
	if err != nil {
		return fmt.Errorf("read file error: %s", srcPath)
	}

	err = os.WriteFile(dstPath, input, 0644)
	if err != nil {
		return fmt.Errorf("write file error: %s", dstPath)
	}

	log.Printf("Скопирован файл: %s -> %s", srcPath, dstPath)

	return nil
}

func WalkAndCopy(srcDir, dstDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %w", err)
		}

		if !info.IsDir() {
			if err := copyFile(path, dstDir); err != nil {
				return fmt.Errorf("copy error: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("walk and move error: %w", err)
	}

	return nil
}
