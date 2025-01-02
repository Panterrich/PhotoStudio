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

func moveFile(srcPath string) error {
	var srcFileName, camera, dstPath string

	srcFileName = filepath.Base(srcPath)
	camera = filepath.Base(filepath.Dir(srcPath))
	fmt.Println(srcPath, srcFileName, camera)

	srcFileName = modifyFileName(srcFileName, camera)
	if srcFileName == "" {
		return fmt.Errorf("empty filename")
	}

	if image.IsRaw(srcFileName) {
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
		return fmt.Errorf("rename file %s -> %s: %w", srcFileName, dstPath, err)
	}

	log.Printf("Скопирован файл: %s -> %s", srcPath, dstPath)

	return nil
}

func MoveImages(srcDir string, nWorkers int) error {
	c, err := collection.NewCollection(srcDir, nWorkers)
	if err != nil {
		return fmt.Errorf("create collection for copying: %w", err)
	}

	jpegSize, rawSize := c.Size()
	images := c.Images()

	p, wg := progressbar.New(2)

	jpegBar := progressbar.Add(p, jpegSize, "Moving JPEGs...")
	rawBar := progressbar.Add(p, rawSize, "Moving RAWs... ")

	result := make(chan error, 2)

	go func() {
		defer wg.Done()

		for jpeg := range images.Jpegs {
			start := time.Now()

			if err := moveFile(jpeg); err != nil {
				result <- fmt.Errorf("moving jpeg file error: %w", err)
				return
			}

			jpegBar.IncrBy(1, time.Since(start))
		}
	}()

	go func() {
		defer wg.Done()

		for raw := range images.Raws {
			start := time.Now()

			if err := moveFile(raw); err != nil {
				result <- fmt.Errorf("moving raw file error: %w", err)
				return
			}

			rawBar.IncrBy(1, time.Since(start))
		}
	}()

	p.Wait()

	for {
		select {
		case err := <-result:
			if err != nil {
				return fmt.Errorf("invalid moving: %w", err)
			}
		default:
			return nil
		}
	}
}
