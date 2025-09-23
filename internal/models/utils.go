package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/arthwr/uma-scrap/internal/config"
)

func ensureDir(dir string) error {
	return os.MkdirAll(dir, os.FileMode(0755))
}

func makeFilename() string {
	now := time.Now().Format(time.DateOnly)
	return fmt.Sprintf("%s%s", now, config.DEF_EVENTS_FILENAME_PATTERN)
}

func writeJSON(path string, data any) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

func cleanupOldFiles(dir, pattern string, keep int) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"+pattern))
	if err != nil {
		return err
	}

	sort.Strings(files)

	if len(files) > keep {
		for _, f := range files[:len(files)-keep] {
			_ = os.Remove(f)
		}
	}

	return nil
}
