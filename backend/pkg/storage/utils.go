// pkg/storage/utils.go
package storage

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// GetFileCategory determines the category of a file based on its extension
func GetFileCategory(filename string) FileCategory {
	ext := strings.ToLower(filepath.Ext(filename))
	
	// Map of extensions to categories
	extensionMap := map[string]FileCategory{
		// Documents
		".pdf":  CategoryDocument,
		".doc":  CategoryDocument,
		".docx": CategoryDocument,
		".txt":  CategoryDocument,
		".rtf":  CategoryDocument,
		".odt":  CategoryDocument,
		".md":   CategoryDocument,
		".csv":  CategoryDocument,
		".xls":  CategoryDocument,
		".xlsx": CategoryDocument,
		".ppt":  CategoryDocument,
		".pptx": CategoryDocument,
		
		// Images
		".jpg":  CategoryImage,
		".jpeg": CategoryImage,
		".png":  CategoryImage,
		".gif":  CategoryImage,
		".webp": CategoryImage,
		".svg":  CategoryImage,
		".bmp":  CategoryImage,
		
		// Audio
		".mp3":  CategoryAudio,
		".wav":  CategoryAudio,
		".ogg":  CategoryAudio,
		".flac": CategoryAudio,
		".m4a":  CategoryAudio,
		
		// Video
		".mp4":  CategoryVideo,
		".mov":  CategoryVideo,
		".avi":  CategoryVideo,
		".mkv":  CategoryVideo,
		".webm": CategoryVideo,
		
		// Archives
		".zip":  CategoryArchive,
		".rar":  CategoryArchive,
		".7z":   CategoryArchive,
		".tar":  CategoryArchive,
		".gz":   CategoryArchive,
	}
	
	if category, ok := extensionMap[ext]; ok {
		return category
	}
	
	return CategoryOther
}

// GenerateObjectName generates a unique object name for storage
func GenerateObjectName(orgID uuid.UUID, category FileCategory, filename string) string {
	// Generate a unique ID to prevent name collisions
	uniqueID := uuid.New().String()
	
	// Extract file extension and base name
	extension := filepath.Ext(filename)
	baseName := filepath.Base(filename[:len(filename)-len(extension)])
	
	// Create a path with organization, category, and unique ID
	return fmt.Sprintf("%s/%s/%s-%s%s", 
		orgID.String(), 
		string(category),
		baseName,
		uniqueID[:8], // Use first 8 characters of UUID for brevity
		extension,
	)
}

// GetSupportedFileTypes returns a map of supported file types and their extensions
func GetSupportedFileTypes() map[string][]string {
	return map[string][]string{
		"documents": {".pdf", ".doc", ".docx", ".txt", ".rtf", ".odt", ".md", ".csv", ".xls", ".xlsx"},
		"images":    {".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp"},
		"audio":     {".mp3", ".wav", ".ogg", ".flac", ".m4a"},
		"video":     {".mp4", ".mov", ".avi", ".mkv", ".webm"},
		"archives":  {".zip", ".rar", ".7z", ".tar", ".gz"},
	}
}