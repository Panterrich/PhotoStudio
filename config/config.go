package config

import (
	"fmt"

	"github.com/Panterrich/PhotoStudio/pkg/image"
	"github.com/spf13/viper"
)

type Config struct {
	viper *viper.Viper

	Cameras image.Cameras
}

func GetConfig(configPath, configName string) (*Config, error) {
	// Initialize viper
	v := viper.New()
	v.AddConfigPath(configPath)
	v.SetConfigName(configName)
	v.SetConfigType("yaml")

	err := v.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("error in reading config: %w", err)
	}

	c := &Config{}
	c.viper = v

	if err = c.setCamerasConfig(); err != nil {
		return nil, fmt.Errorf("failed to set devices config: %w", err)
	}

	return c, nil
}

func (c *Config) setCamerasConfig() error {
	cameras := make(image.Cameras)

	rawBrands := c.viper.GetStringMap("cameras")

	for brand, rawModels := range rawBrands {
		models, ok := rawModels.(map[string]any)
		if !ok {
			return fmt.Errorf("bad models list in config")
		}

		for model, rawPattern := range models {
			pattern, ok := rawPattern.(string)
			if !ok {
				return fmt.Errorf("bad pattern list in config")
			}

			cameras[pattern] = image.Camera{
				Brand: brand,
				Model: model,
			}
		}
	}

	c.Cameras = cameras

	return nil
}
