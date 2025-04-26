package rawmove

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	log "github.com/rs/zerolog/log"

	"github.com/Panterrich/PhotoStudio/pkg/collection"
	"github.com/Panterrich/PhotoStudio/pkg/image"
	"github.com/Panterrich/PhotoStudio/pkg/progressbar"
)

func copyFile(srcPath, dstDir string) error {
	var srcFileName, camera, dstPath string

	srcFileName = filepath.Base(srcPath)
	camera = filepath.Base(dstDir)

	srcFileName = modifyFileName(srcFileName, camera)
	if srcFileName == "" {
		return fmt.Errorf("empty filename: %s", srcFileName)
	}

	if image.IsRaw(srcFileName) {
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

	log.Info().Msgf("Скопирован файл: %s -> %s", srcPath, dstPath)

	return nil
}

func CopyImages(srcDir, dstDir string, nWorkers int) error {
	c, err := collection.NewCollection(srcDir, nWorkers)
	if err != nil {
		return fmt.Errorf("create collection for copying: %w", err)
	}

	jpegSize, rawSize := c.Size()
	images := c.Images()

	p, wg := progressbar.New(2)

	jpegBar := progressbar.Add(p, jpegSize, "Copying JPEGs...")
	rawBar := progressbar.Add(p, rawSize, "Copying RAWs... ")

	result := make(chan error, 2)

	go func() {
		defer wg.Done()

		for jpeg := range images.Jpegs {
			start := time.Now()

			if err := copyFile(jpeg, dstDir); err != nil {
				result <- fmt.Errorf("copying jpeg file error: %w", err)
				return
			}

			jpegBar.IncrBy(1, time.Since(start))
		}

		jpegBar.SetTotal(int64(jpegSize), true)
	}()

	go func() {
		defer wg.Done()

		for raw := range images.Raws {
			start := time.Now()

			if err := copyFile(raw, dstDir); err != nil {
				result <- fmt.Errorf("copying raw file error: %w", err)
				return
			}

			rawBar.IncrBy(1, time.Since(start))
		}

		rawBar.SetTotal(int64(rawSize), true)
	}()

	p.Wait()

	for {
		select {
		case err := <-result:
			if err != nil {
				return fmt.Errorf("invalid copying: %w", err)
			}
		default:
			return nil
		}
	}
}
