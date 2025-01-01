package image

import (
	"errors"
	"fmt"

	"github.com/dsoprea/go-exif"
)

type Camera struct {
	Brand string
	Model string
}

type Cameras map[string]Camera

const ModelTagID = 272

var ErrNotFound = errors.New("camera model not found")

func matchCamera(raw string, cameras Cameras) (string, error) {
	if camera, ok := cameras[raw]; !ok {
		return "", fmt.Errorf("%w: %s", ErrNotFound, raw)
	} else {
		return camera.Model, nil
	}
}

func WhichCamera(path string, cameras Cameras) (string, error) {
	rawExif, err := exif.SearchFileAndExtractExif(path)
	if err != nil {
		return "", fmt.Errorf("can't extract exif from %s: %w", path, err)
	}

	var model string

	im := exif.NewIfdMappingWithStandard()
	ti := exif.NewTagIndex()

	visitor := func(
		_ /* fqIfdPath */ string,
		_ /* ifdIndex */ int,
		tagId uint16,
		_ /* tagType */ exif.TagType,
		valueContext exif.ValueContext) (err error) {

		if tagId != ModelTagID {
			return nil
		}

		model, err = valueContext.FormatFirst()
		if err != nil {
			return fmt.Errorf("can't format value context: %w", err)
		}

		return nil
	}

	_, err = exif.Visit(exif.IfdStandard, im, ti, rawExif, visitor)
	if err != nil {
		return "", fmt.Errorf("can't visit exif: %w", err)
	}

	matchedModel, err := matchCamera(model, cameras)
	if err != nil {
		return "", fmt.Errorf("can't match camera: %w", err)
	}

	return matchedModel, nil
}
