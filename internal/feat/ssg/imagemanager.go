package ssg

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hermesgen/hm"
)

// CLEANUP_ORPHANED_IMAGES controls whether to automatically clean up orphaned images
// Currently set to false - manual cleanup for now, automated later
const CLEANUP_ORPHANED_IMAGES = false

// ImageType represents the type of image being processed
type ImageType string

const (
	ImageTypeContent       ImageType = "content"
	ImageTypeHeader        ImageType = "header"
	ImageTypeSectionHeader ImageType = "section_header"
	ImageTypeBlogHeader    ImageType = "blog_header"
)

// ImageProcessResult contains the result of image processing
type ImageProcessResult struct {
	FilePath     string            // Full file path where image was stored
	RelativePath string            // Relative path for web access
	Filename     string            // Generated filename
	Directory    string            // Directory where image was stored
	Metadata     map[string]string // Image metadata (size, format, etc.)
}

// ImageManager handles all image-related operations
type ImageManager struct {
	hm.Core
	baseImagePath string // Base path for all images (e.g., "./assets/images")
}

// NewImageManager creates a new ImageManager instance

// NewImageManagerWithParams creates an ImageManager with XParams.
func NewImageManager(params hm.XParams) *ImageManager {
	core := hm.NewCore("image-manager", params)
	imagesPath := core.Cfg().StrValOrDef(SSGKey.ImagesPath, "_workspace/documents/assets/images")
	return &ImageManager{
		Core:          core,
		baseImagePath: imagesPath,
	}
}

// ProcessUpload handles the complete upload process for any image type
func (im *ImageManager) ProcessUpload(ctx context.Context, file multipart.File, header *multipart.FileHeader, content *Content, section *Section, imageType ImageType, altText, caption string) (*ImageProcessResult, error) {
	directory, err := im.generateDirectoryPath(content, section, imageType)
	if err != nil {
		return nil, fmt.Errorf("failed to generate directory path: %w", err)
	}

	filename, err := im.generateFilename(content, section, imageType, header.Filename)
	if err != nil {
		return nil, fmt.Errorf("failed to generate filename: %w", err)
	}

	fullDirectory := filepath.Join(im.baseImagePath, directory)
	if err := im.ensureDirectory(fullDirectory); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	if err := im.handleReplacement(fullDirectory, imageType, content, section); err != nil {
		return nil, fmt.Errorf("failed to handle replacement: %w", err)
	}

	fullPath := filepath.Join(fullDirectory, filename)
	if err := im.saveFile(file, fullPath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	metadata := im.extractMetadata(header)

	result := &ImageProcessResult{
		FilePath:     fullPath,
		RelativePath: filepath.Join(directory, filename),
		Filename:     filename,
		Directory:    directory,
		Metadata:     metadata,
	}

	im.Log().Debugf("Upload processed successfully: %s", result.RelativePath)
	return result, nil
}

// generateDirectoryPath creates the directory path based on content hierarchy
func (im *ImageManager) generateDirectoryPath(content *Content, section *Section, imageType ImageType) (string, error) {
	switch imageType {
	case ImageTypeContent, ImageTypeHeader:
		if content == nil {
			return "", fmt.Errorf("content is required for content/header images")
		}

		if section == nil || section.Path == "/" || section.Path == "" {
			return content.Slug(), nil
		}

		return filepath.Join(section.Path, content.Slug()), nil

	case ImageTypeSectionHeader:
		if section == nil {
			return "", fmt.Errorf("section is required for section header images")
		}
		if section.Path == "/" || section.Path == "" {
			return ".", nil
		}
		return section.Path, nil

	case ImageTypeBlogHeader:
		if section == nil {
			return "", fmt.Errorf("section is required for blog header images")
		}
		if section.Path == "/" || section.Path == "" {
			return "blog", nil
		}
		return filepath.Join(section.Path, "blog"), nil

	default:
		return "", fmt.Errorf("unknown image type: %s", imageType)
	}
}

// generateFilename creates the filename based on naming conventions
func (im *ImageManager) generateFilename(content *Content, section *Section, imageType ImageType, originalFilename string) (string, error) {
	timestamp := time.Now().Unix()
	extension := filepath.Ext(originalFilename)

	switch imageType {
	case ImageTypeContent:
		if content == nil {
			return "", fmt.Errorf("content is required for content images")
		}
		return fmt.Sprintf("%s_%d%s", content.Slug(), timestamp, extension), nil

	case ImageTypeHeader:
		if content == nil {
			return "", fmt.Errorf("content is required for header images")
		}
		return fmt.Sprintf("%s_header_%d%s", content.Slug(), timestamp, extension), nil

	case ImageTypeSectionHeader:
		if section == nil {
			return "", fmt.Errorf("section is required for section header images")
		}
		return fmt.Sprintf("section_header_%d%s", timestamp, extension), nil

	case ImageTypeBlogHeader:
		if section == nil {
			return "", fmt.Errorf("section is required for blog header images")
		}
		return fmt.Sprintf("blog_header_%d%s", timestamp, extension), nil

	default:
		return "", fmt.Errorf("unknown image type: %s", imageType)
	}
}

// ensureDirectory creates directory if it doesn't exist
func (im *ImageManager) ensureDirectory(path string) error {
	return os.MkdirAll(path, 0755)
}

// handleReplacement removes old files for single-image types
func (im *ImageManager) handleReplacement(directory string, imageType ImageType, content *Content, section *Section) error {
	if imageType != ImageTypeHeader && imageType != ImageTypeSectionHeader && imageType != ImageTypeBlogHeader {
		return nil
	}

	var pattern string
	switch imageType {
	case ImageTypeHeader:
		if content == nil {
			return fmt.Errorf("content is required for header image replacement")
		}
		pattern = content.Slug() + "_header_*"
	case ImageTypeSectionHeader:
		pattern = "section_header_*"
	case ImageTypeBlogHeader:
		pattern = "blog_header_*"
	}

	files, err := filepath.Glob(filepath.Join(directory, pattern))
	if err != nil {
		return fmt.Errorf("failed to find existing files: %w", err)
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil {
			im.Log().Errorf("Failed to remove old file %s: %v", file, err)
		}
	}

	return nil
}

func (im *ImageManager) saveFile(src multipart.File, destPath string) error {
	if _, err := src.Seek(0, 0); err != nil {
		return fmt.Errorf("failed to reset file pointer: %w", err)
	}

	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to copy file content: %w", err)
	}

	return nil
}

// extractMetadata extracts basic metadata from the uploaded file
func (im *ImageManager) extractMetadata(header *multipart.FileHeader) map[string]string {
	metadata := make(map[string]string)

	metadata["original_filename"] = header.Filename
	metadata["content_type"] = header.Header.Get("Content-Type")
	metadata["size"] = fmt.Sprintf("%d", header.Size)
	metadata["upload_time"] = time.Now().Format(time.RFC3339)

	return metadata
}

// GetContentImages returns all images for a specific content
func (im *ImageManager) GetContentImages(ctx context.Context, content *Content, section *Section) ([]string, error) {
	directory, err := im.generateDirectoryPath(content, section, ImageTypeContent)
	if err != nil {
		return nil, err
	}

	fullDirectory := filepath.Join(im.baseImagePath, directory)
	if _, err := os.Stat(fullDirectory); os.IsNotExist(err) {
		return []string{}, nil // No images directory yet
	}

	files, err := filepath.Glob(filepath.Join(fullDirectory, "*"))
	if err != nil {
		return nil, fmt.Errorf("failed to list files: %w", err)
	}

	var images []string
	for _, file := range files {
		filename := filepath.Base(file)
		if !strings.Contains(filename, "_header_") {
			images = append(images, filepath.Join(directory, filename))
		}
	}

	return images, nil
}

// CleanupOrphanedImages removes images that no longer have associated content
// This is a placeholder for future implementation
func (im *ImageManager) CleanupOrphanedImages(ctx context.Context) error {
	if !CLEANUP_ORPHANED_IMAGES {
		im.Log().Debug("Orphaned image cleanup is disabled")
		return nil
	}

	// TODO: Implement orphaned image detection and cleanup
	im.Log().Debug("Orphaned image cleanup not yet implemented")
	return nil
}

// DeleteImage deletes an image file by its relative path
func (im *ImageManager) DeleteImage(ctx context.Context, relativePath string) error {
	if relativePath == "" {
		return nil // Nothing to delete
	}

	fullPath := filepath.Join(im.baseImagePath, relativePath)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return nil // File doesn't exist, consider it deleted
	}

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete image file %s: %w", fullPath, err)
	}

	return nil
}
