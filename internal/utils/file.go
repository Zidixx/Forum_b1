package utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const MaxUploadSize = 20 << 20 // 20 MB

var allowedTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
}

var allowedExts = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
}

func SaveUploadedFile(file multipart.File, header *multipart.FileHeader, uploadDir string) (string, error) {
	if header.Size > MaxUploadSize {
		return "", fmt.Errorf("le fichier est trop volumineux (max 20 Mo)")
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !allowedExts[ext] {
		return "", fmt.Errorf("format de fichier non autorisé (JPEG, PNG, GIF uniquement)")
	}

	buf := make([]byte, 512)
	if _, err := file.Read(buf); err != nil {
		return "", fmt.Errorf("impossible de lire le fichier")
	}
	contentType := http.DetectContentType(buf)
	if !allowedTypes[contentType] {
		return "", fmt.Errorf("type de fichier invalide")
	}
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return "", fmt.Errorf("erreur lecture fichier")
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("impossible de créer le dossier uploads")
	}

	filename := GenerateUUID() + ext
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("impossible de sauvegarder le fichier")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("erreur lors de la sauvegarde")
	}

	return filename, nil
}
