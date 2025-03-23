package storage

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type File struct{}

func NewFileStorage() *File {
	return &File{}
}

func (s *File) Upload(file io.Reader, fileName string) error {
	var buff bytes.Buffer
	_, err := io.Copy(&buff, file)
	if err != nil {
		return fmt.Errorf("couldn't copy file to buffer: %v", err)
	}
	filePath := fmt.Sprintf("%s/%s", viper.GetString("STORAGE_FILE_PATH"), fileName)
	err = os.WriteFile(filePath, buff.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("couldn't write file to disk: %v", err)
	}
	return nil
}

func (s *File) ListFiles(folder string) ([]StoredFile, error) {
	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, fmt.Errorf("couldn't list files in folder: %v", err)
	}

	if len(entries) == 0 {
		return nil, nil
	}

	files := make([]StoredFile, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		timestamp, err := parseTimeFromFilename(entry.Name())
		if err != nil {
			continue
		}

		files = append(files, StoredFile{
			Name:      entry.Name(),
			Timestamp: timestamp,
		})
	}

	return files, nil
}

func (s *File) Delete(filename string) error {
	fullPath := filepath.Join(viper.GetString("STORAGE_FILE_PATH"), filename)
	err := os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("couldn't delete file: %v", err)
	}
	return nil
}
