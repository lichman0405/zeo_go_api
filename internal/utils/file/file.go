package file

import (
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var validExtensions = map[string]bool{
	".cif":     true,
	".cssr":    true,
	".v1":      true,
	".arc":     true,
	".cif.gz":  true,
	".cssr.gz": true,
}

func IsValidStructureFile(filename string) bool {
	filename = strings.ToLower(filename)
	for ext := range validExtensions {
		if strings.HasSuffix(filename, ext) {
			return true
		}
	}
	return false
}

func SaveUploadedFile(file *multipart.FileHeader, prefix string) (string, error) {
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Sanitize filename
	safeName := sanitizeFilename(file.Filename)

	// Generate unique filename
	uniqueID := fmt.Sprintf("%s_%d_%d", prefix, time.Now().UnixNano(), rand.Intn(10000))
	ext := filepath.Ext(safeName)
	if strings.HasSuffix(strings.ToLower(safeName), ".gz") {
		// Handle gzipped files properly
		base := strings.TrimSuffix(safeName, ".gz")
		ext = filepath.Ext(base) + ".gz"
	}

	// Ensure workspace directory exists
	workspace := "./workspace"
	if err := os.MkdirAll(workspace, 0700); err != nil {
		return "", err
	}

	// Ensure path is within workspace
	filename := fmt.Sprintf("%s%s", uniqueID, ext)
	fullPath := filepath.Join(workspace, filename)

	// Validate final path
	if !strings.HasPrefix(fullPath, workspace) {
		return "", fmt.Errorf("invalid file path")
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		_ = os.Remove(fullPath) // Clean up on error
		return "", err
	}

	return fullPath, nil
}

func sanitizeFilename(filename string) string {
	// Remove path separators and sanitize
	filename = filepath.Base(filename)
	filename = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '.' || r == '-' || r == '_' {
			return r
		}
		return '_'
	}, filename)

	// Ensure filename is not empty
	if filename == "" || filename == "." {
		return "uploaded_file"
	}
	return filename
}

func GenerateFileHash(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func CleanupFile(filePath string) {
	if filePath != "" {
		os.Remove(filePath)
	}
}

func CleanupDirectory(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && time.Since(info.ModTime()) > 24*time.Hour {
			return os.Remove(path)
		}
		return nil
	})
}

func GetFileContent(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}
