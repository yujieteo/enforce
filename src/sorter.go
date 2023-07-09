package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileSorter represents the template for sorting files.
type FileSorter struct {
	FolderPath string
}

// Execute executes the template for sorting files.
func (s *FileSorter) Execute() error {
	err := filepath.Walk(s.FolderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		extension := strings.ToLower(filepath.Ext(path))
		destFolder := ""

		switch extension {
		case ".pdf", ".djvu", ".epub", ".html", ".docx", ".md", ".tex", ".txt", ".doc", ".pptx", ".ipynb":
			destFolder = filepath.Join("doc", strings.TrimSuffix(filepath.Base(path), extension))
		case ".rst", ".rth", ".cdb", ".ls-dyna", ".db", ".dbb", ".esav", ".out":
			destFolder = "job"
		case ".mkv", ".mp4", ".aac", ".flac", ".wav", ".avi", ".png", ".jpeg", ".mov", ".wmv", ".jpg", ".mp3":
			destFolder = filepath.Join("media", strings.TrimSuffix(filepath.Base(path), extension))
		case ".py", ".go", ".ans", ".inp", ".c", ".m", ".for", ".cpp", ".java", ".scala", ".php", ".sh", ".asm", ".h", ".dat":
			destFolder = filepath.Join("src", strings.TrimSuffix(filepath.Base(path), extension))
		case ".exe":
			destFolder = "bin"
		default:
			destFolder = "data"
		}

		destFolderPath := filepath.Join(s.FolderPath, destFolder)
		err = os.MkdirAll(destFolderPath, 0755)
		if err != nil {
			return err
		}

		destFilePath := filepath.Join(destFolderPath, filepath.Base(path))
		err = os.Rename(path, destFilePath)
		if err != nil {
			return err
		}

		fmt.Printf("Moved '%s' to '%s'\n", path, destFilePath)
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to sort files: %w", err)
	}

	return nil
}
