package main

import (
	"os"
)

// Directory represents a directory in the file system.
type Directory struct {
	path       string
	operations []FileOperation
}

// AddOperation adds a file operation to the directory.
func (d *Directory) AddOperation(op FileOperation) {
	d.operations = append(d.operations, op)
}

// ExecuteOperations executes all file operations in the directory.
func (d *Directory) ExecuteOperations() error {
	for _, op := range d.operations {
		err := op.Execute()
		if err != nil {
			return err
		}
	}
	return nil
}

// RecursiveDirectory represents a directory with recursive operations.
type RecursiveDirectory struct {
	*Directory
	subdirectories []*RecursiveDirectory
}

// AddSubdirectory adds a subdirectory to the recursive directory.
func (r *RecursiveDirectory) AddSubdirectory(dir *RecursiveDirectory) {
	r.subdirectories = append(r.subdirectories, dir)
}

// ExecuteOperations executes all file operations in the recursive directory and its subdirectories.
func (r *RecursiveDirectory) ExecuteOperations() error {
	err := r.Directory.ExecuteOperations()
	if err != nil {
		return err
	}
	for _, subdir := range r.subdirectories {
		err := subdir.ExecuteOperations()
		if err != nil {
			return err
		}
	}
	return nil
}

// Helper function to check if a directory is empty
func isDirectoryEmpty(dirPath string) (bool, error) {
	f, err := os.Open(dirPath)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == nil {
		// Directory is not empty
		return false, nil
	}

	if err == os.ErrNotExist {
		// Directory doesn't exist
		return false, err
	}

	if err == os.ErrPermission {
		// Permission denied
		return false, err
	}

	// Directory is empty
	return true, nil
}
