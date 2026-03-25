package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func MigrateStorage(targetPath string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	oldPath := filepath.Join(home, ".pith")

	// If old path doesn't exist, nothing to migrate
	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return nil
	}

	// If target is the same as old, nothing to do
	if oldPath == targetPath {
		return nil
	}

	// Create target directory
	if err := os.MkdirAll(targetPath, 0755); err != nil {
		return err
	}

	files := []string{"pith.db", "config.json"}
	for _, f := range files {
		src := filepath.Join(oldPath, f)
		dst := filepath.Join(targetPath, f)

		// Skip if source doesn't exist
		if _, err := os.Stat(src); os.IsNotExist(err) {
			continue
		}

		// Skip if destination already exists (don't overwrite)
		if _, err := os.Stat(dst); err == nil {
			continue
		}

		fmt.Printf("[Pith] Migrating %s to %s...\n", f, targetPath)
		if err := copyFile(src, dst); err != nil {
			return fmt.Errorf("failed to migrate %s: %w", f, err)
		}
		
		// Optional: rename old file to .bak instead of deleting immediately for safety
		_ = os.Rename(src, src+".bak")
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
