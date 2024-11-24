package main

import (
	"fmt"
	"os"
	"path"

	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	rawmove "github.com/Panterrich/PhotoStudio/pkg/raw-move"
)

type Config struct {
	srcDir string
	dstDir string
}

var (
	cfg Config

	root = &cobra.Command{
		Use:   "photostudio",
		Short: "Server for storing metrics",
		Long:  "Server for storing metrics",
	}

	moveCommand = &cobra.Command{
		Use:  "move",
		RunE: move,
	}

	removeCommand = &cobra.Command{
		Use:  "remove",
		RunE: remove,
	}
)

func init() {
	moveCommand.PersistentFlags().StringVarP(&cfg.srcDir, "input", "i", "", "input dir for src")
	moveCommand.PersistentFlags().StringVarP(&cfg.dstDir, "output", "o", "", "output dir for src and raw")
	moveCommand.MarkPersistentFlagRequired("input")
	moveCommand.MarkPersistentFlagRequired("output")

	removeCommand.PersistentFlags().StringVarP(&cfg.srcDir, "input", "i", "", "input dir for src")
	removeCommand.MarkPersistentFlagRequired("input")

	root.AddCommand(moveCommand, removeCommand)
}

func move(_ *cobra.Command, _ []string) error {
	srcDir := cfg.srcDir
	dstDir := cfg.dstDir

	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dstDir, 0755); err != nil {
			log.Error().Msgf("Не удалось создать директорию назначения: %v", err)
		}
	}

	rawDir := path.Join(dstDir, rawmove.RawDir)

	if _, err := os.Stat(rawDir); os.IsNotExist(err) {
		if err := os.Mkdir(rawDir, 0755); err != nil {
			log.Error().Msgf("Не удалось создать RAW директорию: %v", err)
		}
	}

	if err := rawmove.WalkAndMove(srcDir, dstDir); err != nil {
		log.Error().Msgf("Ошибка при копировании: %v", err)
	}

	fmt.Println("Копирование завершено успешно!")

	return nil
}

func remove(_ *cobra.Command, _ []string) error {
	if err := rawmove.RemoveUnnecessaryRaws(cfg.srcDir); err != nil {
		return fmt.Errorf("remove command: %w", err)
	}

	fmt.Println("Удаление завершено")

	return nil
}

func main() {
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}
