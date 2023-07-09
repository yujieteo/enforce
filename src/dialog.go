package main

import (
	"fmt"

	"github.com/sqweek/dialog"
)

// DirectoryDialogFactory is a factory that creates a directory dialog.
type DirectoryDialogFactory struct{}

// Dialog is an interface representing a directory dialog.
type Dialog interface {
	Browse() (string, error)
}

// CreateDialog creates a new directory dialog.
func (f *DirectoryDialogFactory) CreateDialog() Dialog {
	return &DirectoryDialog{}
}

// DirectoryDialog is a directory dialog implementation.
type DirectoryDialog struct{}

// Browse displays the directory dialog and returns the selected path.
func (d *DirectoryDialog) Browse() (string, error) {
	projectPath, err := dialog.Directory().Title("Select project directory").Browse()
	if err != nil {
		return "", fmt.Errorf("failed to select project directory: %w", err)
	}
	return projectPath, nil
}
