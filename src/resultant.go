package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// TextFileFactory is a factory that creates various types of text files.
type TextFileFactory struct {
	ProjectPath string
}

// CreateGitignore creates a .gitignore file in the project path.
func (f *TextFileFactory) CreateGitignore() error {
	gitignorePath := filepath.Join(f.ProjectPath, ".gitignore")
	if _, err := os.Stat(gitignorePath); !os.IsNotExist(err) {
		return fmt.Errorf(".gitignore already exists in the project path")
	}

	gitignoreContent := []byte("# Generated .gitignore file\n\n# Ignore build and temporary files\nbuild/\n*.tmp\n\n# Ignore IDE-specific files\n.idea/\n.vscode/\n\n# Ignore compiled binaries\n*.exe\n*.dll\n\n# Ignore system and OS files\nThumbs.db\n.DS_Store\n")

	err := ioutil.WriteFile(gitignorePath, gitignoreContent, 0644)
	if err != nil {
		return fmt.Errorf("failed to create .gitignore file: %w", err)
	}

	return nil
}
