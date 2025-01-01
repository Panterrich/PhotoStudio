package rawmove

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/Panterrich/PhotoStudio/pkg/image"
)

func RemoveUnnecessaryRaws(srcDirPath string) error {
	rawDirPath := path.Join(srcDirPath, RawDir)

	rawDir, err := os.Open(rawDirPath)
	if err != nil {
		return fmt.Errorf("raw dir open: %w", err)
	}
	defer rawDir.Close()

	raws, err := rawDir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("raw dir read: %w", err)
	}

	srcDir, err := os.Open(srcDirPath)
	if err != nil {
		return fmt.Errorf("src dir open: %w", err)
	}
	defer srcDir.Close()

	src, err := srcDir.Readdir(-1)
	if err != nil {
		return fmt.Errorf("src dir read: %w", err)
	}

	srcSet := make(map[string]struct{})
	for _, file := range src {
		srcSet[strings.Split(file.Name(), ".")[0]] = struct{}{}
	}

	for _, raw := range raws {
		if !image.IsRaw(raw.Name()) {
			continue
		}

		if _, ok := srcSet[strings.Split(raw.Name(), ".")[0]]; !ok {
			filePath := path.Join(rawDirPath, raw.Name())

			if err := os.Remove(filePath); err != nil {
				return fmt.Errorf("can't remove %s: %w", filePath, err)
			}
		}
	}

	return nil
}
