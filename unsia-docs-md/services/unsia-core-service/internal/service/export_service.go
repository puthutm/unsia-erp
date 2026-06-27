package service

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

// FileService handles file storage
type FileService struct {
	uploadDir string
}

func NewFileService(uploadDir string) *FileService {
	return &FileService{uploadDir: uploadDir}
}

type UploadedFile struct {
	ID          string    `json:"id"`
	Filename   string    `json:"filename"`
	OriginalName string  `json:"original_name"`
	Size       int64     `json:"size"`
	ContentType string   `json:"content_type"`
	Path       string    `json:"path"`
	URL        string    `json:"url"`
	UploadedAt time.Time `json:"uploaded_at"`
}

// Allowed file types
var allowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

var allowedDocumentTypes = map[string]bool{
	"application/pdf":                       true,
	"application/msword":                    true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/vnd.ms-excel":              true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":      true,
	"text/csv":                             true,
}

// Max file sizes (in bytes)
const (
	MaxImageSize    int64 = 10 * 1024 * 1024  // 10MB
	MaxDocumentSize int64 = 50 * 1024 * 1024 // 50MB
)

// ValidateFileType validates file type
func (s *FileService) ValidateFileType(contentType string, category string) error {
	switch category {
	case "image":
		if !allowedImageTypes[contentType] {
			return fmt.Errorf("invalid image type: %s", contentType)
		}
	case "document":
		if !allowedDocumentTypes[contentType] {
			return fmt.Errorf("invalid document type: %s", contentType)
		}
	default:
		if !allowedImageTypes[contentType] && !allowedDocumentTypes[contentType] {
			return fmt.Errorf("invalid file type: %s", contentType)
		}
	}
	return nil
}

// ValidateFileSize validates file size
func (s *FileService) ValidateFileSize(size int64, category string) error {
	maxSize := MaxDocumentSize
	if category == "image" {
		maxSize = MaxImageSize
	}

	if size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed: %d > %d", size, maxSize)
	}
	return nil
}

// GenerateFilePath generates file path
func (s *FileService) GenerateFilePath(filename string) string {
	ext := filepath.Ext(filename)
	dateDir := time.Now().Format("2006/01")
	newFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	return filepath.Join(dateDir, newFilename)
}

// SaveFile saves file to disk
func (s *FileService) SaveFile(data []byte, relativePath string) error {
	fullPath := filepath.Join(s.uploadDir, relativePath)
	
	// Create directory if not exists
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write file
	return os.WriteFile(fullPath, data, 0644)
}

// DeleteFile deletes file from disk
func (s *FileService) DeleteFile(relativePath string) error {
	fullPath := filepath.Join(s.uploadDir, relativePath)
	return os.Remove(fullPath)
}

// GenerateThumbnail generates thumbnail for image (simplified)
// In production, use proper image processing library
func (s *FileService) GenerateThumbnail(imagePath string) error {
	// This would use something like bimg or imaging library
	return nil
}

// FileToJSON converts file to JSON for storage
func (s *FileService) FileToJSON(file *UploadedFile) (string, error) {
	data, err := json.Marshal(file)
	return string(data), err
}

// JSONToFile converts JSON to file struct
func (s *FileService) JSONToFile(jsonStr string) (*UploadedFile, error) {
	var file UploadedFile
	err := json.Unmarshal([]byte(jsonStr), &file)
	return &file, err
}

// ProfileImageService handles profile images
type ProfileImageService struct {
	FileService *FileService
}

func NewProfileImageService(uploadDir string) *ProfileImageService {
	return &ProfileImageService{
		FileService: NewFileService(uploadDir),
	}
}

// UploadProfileImage uploads profile image
func (s *ProfileImageService) UploadProfileImage(userID string, data []byte, contentType string) (*UploadedFile, error) {
	// Validate
	if err := s.FileService.ValidateFileType(contentType, "image"); err != nil {
		return nil, err
	}

	if err := s.FileService.ValidateFileSize(int64(len(data)), "image"); err != nil {
		return nil, err
	}

	// Generate filename
	filename := fmt.Sprintf("%s_profile.%s", userID, getExtension(contentType))
	path := s.FileService.GenerateFilePath(filename)
	path = filepath.Join("profile", path)

	// Save
	if err := s.FileService.SaveFile(data, path); err != nil {
		return nil, err
	}

	return &UploadedFile{
		ID:           uuid.New().String(),
		Filename:     filename,
		OriginalName: filename,
		Size:        int64(len(data)),
		ContentType: contentType,
		Path:        path,
		UploadedAt:  time.Now(),
	}, nil
}

// DeleteProfileImage deletes profile image
func (s *ProfileImageService) DeleteProfileImage(userID, imagePath string) error {
	// Only delete if it's actually the user's current profile image
	return s.FileService.DeleteFile(imagePath)
}

func getExtension(contentType string) string {
	extensions := map[string]string{
		"image/jpeg": "jpg",
		"image/png":  "png",
		"image/gif": "gif",
		"image/webp": "webp",
	}
	if ext, ok := extensions[contentType]; ok {
		return ext
	}
	return "jpg"
}
