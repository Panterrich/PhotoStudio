package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"
	log "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	rawmove "github.com/Panterrich/PhotoStudio/pkg/raw-move"
)

type Flags struct {
	srcDir string
	dstDir string

	debug bool

	nWorkers int
}

var (
	cfg Flags

	root = &cobra.Command{
		Use:              "photostudio",
		Short:            "Script for images processing",
		Long:             "Script for images processing",
		PersistentPreRun: initCfg,
	}

	copyCommand = &cobra.Command{
		Use:  "copy",
		RunE: cmdCopy,
	}

	moveCommand = &cobra.Command{
		Use:    "move",
		Hidden: true,
		RunE:   cmdMove,
	}

	removeCommand = &cobra.Command{
		Use:  "remove",
		RunE: cmdRemove,
	}
)

func init() {
	copyCommand.PersistentFlags().StringVarP(&cfg.srcDir, "input", "i", ".", "input dir for src")
	copyCommand.PersistentFlags().StringVarP(&cfg.dstDir, "output", "o", "", "output dir for src and raw")
	copyCommand.PersistentFlags().IntVarP(&cfg.nWorkers, "jobs", "j", 1, "nWorkers")
	copyCommand.MarkPersistentFlagRequired("output")

	moveCommand.PersistentFlags().StringVarP(&cfg.srcDir, "input", "i", ".", "input dir for src")
	moveCommand.PersistentFlags().IntVarP(&cfg.nWorkers, "jobs", "j", 1, "nWorkers")

	removeCommand.PersistentFlags().StringVarP(&cfg.srcDir, "input", "i", ".", "input dir for src")

	root.AddCommand(copyCommand, moveCommand, removeCommand)

	root.PersistentFlags().BoolVarP(&cfg.debug, "debug", "d", false, "enable debug logs")
}

func initCfg(_ *cobra.Command, _ []string) {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stdout})

	if cfg.debug {
		logger = logger.Level(zerolog.DebugLevel)
	} else {
		logger = logger.Level(zerolog.ErrorLevel)
	}

	log.Logger = logger
}

func cmdCopy(_ *cobra.Command, _ []string) error {
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

	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		log.Error().Msgf("Не удалось получить абсолютный путь: %v", err)
	}

	dstDir, err = filepath.Abs(dstDir)
	if err != nil {
		log.Error().Msgf("Не удалось получить абсолютный путь: %v", err)
	}

	if err = rawmove.CopyImages(srcDir, dstDir, cfg.nWorkers); err != nil {
		log.Error().Msgf("Ошибка при копировании: %v", err)
	}

	fmt.Println("Копирование завершено успешно!")

	return nil
}

func cmdMove(_ *cobra.Command, _ []string) error {
	srcDir := cfg.srcDir

	rawDir := path.Join(srcDir, rawmove.RawDir)

	if _, err := os.Stat(rawDir); os.IsNotExist(err) {
		if err := os.Mkdir(rawDir, 0755); err != nil {
			log.Error().Msgf("Не удалось создать RAW директорию: %v", err)
		}
	}

	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		log.Error().Msgf("Не удалось получить абсолютный путь: %v", err)
	}

	if err = rawmove.MoveImages(srcDir, cfg.nWorkers); err != nil {
		log.Error().Msgf("Ошибка при перемещении: %v", err)
	}

	fmt.Println("Копирование завершено успешно!")

	return nil
}

func cmdRemove(_ *cobra.Command, _ []string) error {
	srcDir, err := filepath.Abs(cfg.srcDir)
	if err != nil {
		log.Error().Msgf("Не удалось получить абсолютный путь: %v", err)
	}

	fmt.Printf("Удалить лишние RAW из %s? [y/N]\n", srcDir)

	var s string
	_, _ = fmt.Scanf("%s", &s)

	if strings.ToLower(s) != "y" {
		fmt.Println("Удаление отменено")
		return nil
	}

	if err := rawmove.RemoveUnnecessaryRaws(srcDir); err != nil {
		return fmt.Errorf("remove command: %w", err)
	}

	fmt.Println("Удаление завершено")

	return nil
}

func main() {
	if err := root.Execute(); err != nil {
		return
	}
}
