package rawmove

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/rs/zerolog/log"
)

const (
	RawDir = "RAW"
)

func isRawFile(fileName string) bool {
	rawFileExtension := map[string]struct{}{
		".cr2": {},
		".nef": {},
	}

	fileName = strings.ToLower(fileName)

	for extension := range rawFileExtension {
		if strings.HasSuffix(fileName, extension) {
			return true
		}
	}

	return false
}

func modifyFileName(fileName string, camera string) string {
	s := strings.FieldsFunc(fileName, func(r rune) bool {
		return r == '_' || r == ' ' || r == '(' || r == ')'
	})

	if len(s) == 0 {
		return ""
	}

	s[0] = camera

	return strings.Join(s, "_")
}

func moveFile(srcPath, dstDir string) error {
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

func WalkAndMove(srcDir, dstDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %w", err)
		}

		if !info.IsDir() {
			if err := moveFile(path, dstDir); err != nil {
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
