package rawmove

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/rs/zerolog/log"
)

func moveFile(srcPath string) error {
	var srcFileName, camera, dstPath string

	srcFileName = filepath.Base(srcPath)
	camera = filepath.Base(filepath.Dir(srcPath))
	fmt.Println(srcPath, srcFileName, camera)

	srcFileName = modifyFileName(srcFileName, camera)
	if srcFileName == "" {
		return fmt.Errorf("empty filename")
	}

	if isRawFile(srcFileName) {
		dstPath = filepath.Join(camera, RawDir, srcFileName)
	} else {
		dstPath = filepath.Join(camera, srcFileName)
	}

	if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
		log.Warn().Msgf("файл %s уже существует", dstPath)
		return nil
	}

	err := os.Rename(srcFileName, dstPath)
	if err != nil {
		return fmt.Errorf("rename file %s -> %s: %v", srcFileName, dstPath, err)
	}

	log.Printf("Скопирован файл: %s -> %s", srcPath, dstPath)

	return nil
}

func WalkAndMove(srcDir string) error {
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("walk error: %w", err)
		}

		if !info.IsDir() {
			if err := moveFile(path); err != nil {
				return fmt.Errorf("move error: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("walk and move error: %w", err)
	}

	return nil
}
