package utils

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

const (
	FOLDER_PERMISSION = 0755
)

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()
	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// CopyDir recursively copies a directory from src to dst
func CopyDir(src, dst string) error {

	// Create the destination directory with the same permissions
	err := os.MkdirAll(dst, FOLDER_PERMISSION)
	if err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err // Propagate errors
		}

		// Find the relative path of the item from the source directory
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		// Create the full destination path
		destPath := filepath.Join(dst, relPath)

		if d.IsDir() {
			// Create directory at destination
			return os.MkdirAll(destPath, FOLDER_PERMISSION)
		} else {
			// Copy file
			return CopyFile(path, destPath)
		}
	})
}
