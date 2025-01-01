package image_test

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charlievieth/fastwalk"
	"github.com/stretchr/testify/assert"

	"github.com/Panterrich/PhotoStudio/config"
	"github.com/Panterrich/PhotoStudio/pkg/image"
)

func TestWhichCamera(t *testing.T) {
	testDir := "testdata"

	cfg, err := config.GetConfig("../../config", "mipt-photo.yaml")
	assert.NoError(t, err)

	cameras := cfg.Cameras

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsRegular() {
			return nil
		}

		assert.NoError(t, err)

		var camera string

		camera, err = image.WhichCamera(path, cameras)

		assert.NoError(t, err)
		assert.True(t, strings.HasPrefix(filepath.Base(path), strings.ToUpper(camera)))

		return nil
	}

	assert.NoError(t, fastwalk.Walk(&fastwalk.DefaultConfig, testDir, walkFn))
}
