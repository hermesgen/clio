package ssg

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

func CopyStaticAssets(assetsFS embed.FS, targetDir string) error {

	return fs.WalkDir(assetsFS, "assets/ssg/static", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking directory: %w", err)
		}

		destPath := filepath.Join(targetDir, path[len("assets/ssg"):])

		if d.IsDir() {
			if err := os.MkdirAll(destPath, 0755); err != nil {
				return fmt.Errorf("cannot create directory: %w", err)
			}
			return nil
		}

		return copyFile(assetsFS, path, destPath)
	})
}

// CopyDynamicImages copies all dynamic images from assets/images to html/static/images
func CopyDynamicImages(sourceDir, targetDir string) error {
	sourceImagesDir := filepath.Join(sourceDir, "assets", "images")
	targetImagesDir := filepath.Join(targetDir, "static", "images")

	if _, err := os.Stat(sourceImagesDir); os.IsNotExist(err) {
		return nil
	}

	if err := os.MkdirAll(targetImagesDir, 0755); err != nil {
		return fmt.Errorf("cannot create images directory: %w", err)
	}

	return filepath.Walk(sourceImagesDir, func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking source directory: %w", err)
		}

		relPath, err := filepath.Rel(sourceImagesDir, srcPath)
		if err != nil {
			return fmt.Errorf("cannot get relative path: %w", err)
		}

		dstPath := filepath.Join(targetImagesDir, relPath)

		if info.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return fmt.Errorf("cannot create directory: %w", err)
			}
			return nil
		}

		return copyFileFromFS(srcPath, dstPath)
	})
}

func copyFile(assetsFS embed.FS, srcPath, dstPath string) error {
	srcFile, err := assetsFS.Open(srcPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("cannot copy file: %w", err)
	}

	return nil
}

// copyFileFromFS copies a file from filesystem to filesystem
func copyFileFromFS(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("cannot open source file: %w", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("cannot create destination file: %w", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("cannot copy file: %w", err)
	}

	return nil
}
