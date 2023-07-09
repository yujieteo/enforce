package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// FileOperation represents a file operation.
type FileOperation interface {
	Execute() error
}

// MoveFileOperation represents a move file operation.
type MoveFileOperation struct {
	sourcePath string
	destPath   string
}

// Execute executes the move file operation.
func (m *MoveFileOperation) Execute() error {
	err := os.Rename(m.sourcePath, m.destPath)
	if err != nil {
		return fmt.Errorf("failed to move file '%s' to '%s': %w", m.sourcePath, m.destPath, err)
	}
	return nil
}

// RenameFileOperation represents a rename file operation.
type RenameFileOperation struct {
	filePath string
	newName  string
}

// Execute executes the rename file operation.
func (r *RenameFileOperation) Execute() error {
	oldFilePath := r.filePath
	newFileName := transformFileName(filepath.Base(oldFilePath))
	newFilePath := filepath.Join(filepath.Dir(oldFilePath), newFileName)

	if oldFilePath != newFilePath {
		err := os.Rename(oldFilePath, newFilePath)
		if err != nil {
			return fmt.Errorf("failed to rename file '%s' to '%s': %w", oldFilePath, newFilePath, err)
		}
	}
	return nil
}

func transformFileName(fileName string) string {
	fileName = regexp.MustCompile(`[\s-]`).ReplaceAllString(fileName, "_")
	fileName = strings.ToLower(fileName)
	fileName = regexp.MustCompile(`_+`).ReplaceAllString(fileName, "_")
	return fileName
}

// CreateDirectoryOperation represents a create directory operation.
type CreateDirectoryOperation struct {
	dirPath string
}

// Execute executes the create directory operation.
func (c *CreateDirectoryOperation) Execute() error {
	err := os.MkdirAll(c.dirPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create directory '%s': %w", c.dirPath, err)
	}
	return nil
}

// RemoveDirectoryOperation represents a remove directory operation.
type RemoveDirectoryOperation struct {
	dirPath string
}

// Execute executes the remove directory operation.
func (r *RemoveDirectoryOperation) Execute() error {
	err := os.Remove(r.dirPath)
	if err != nil {
		return fmt.Errorf("failed to remove directory '%s': %w", r.dirPath, err)
	}
	return nil
}
