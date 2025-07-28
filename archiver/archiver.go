package archiver

import (
	"archive/zip"
	"errors"
	"fmt"
	"image_service/internal/model"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const archivesDir = "data/archives"

type ArchiverService struct{}

func NewArchiverService() *ArchiverService {
	_ = os.MkdirAll(archivesDir, 0755)
	return &ArchiverService{}
}

func (a *ArchiverService) Archive(task *model.Task) {
	fmt.Printf("Start archive %s\n", task.ID)

	filesDir := filepath.Join(archivesDir, task.ID)
	if err := os.MkdirAll(filesDir, 0755); err != nil {
		fmt.Printf("error %v\n", err)
		task.Status = model.StatusFailed
		return
	}

	var localFiles []string
	var failed []string

	for _, url := range task.Files {
		filePath, err := a.downloadFile(url, filesDir)
		if err != nil {
			fmt.Printf("error download %s: %v\n", url, err)
			failed = append(failed, url)
			continue
		}
		localFiles = append(localFiles, filePath)
	}

	if len(localFiles) == 0 {
		task.Status = model.StatusFailed
		task.FailedFiles = failed
		return
	}

	archivePath := filepath.Join(archivesDir, fmt.Sprintf("%s.zip", task.ID))
	if err := a.createZip(archivePath, localFiles); err != nil {
		fmt.Printf("error archive %v\n", err)
		task.Status = model.StatusFailed
		task.FailedFiles = failed
		return
	}

	task.Status = model.StatusComplete
	task.URL = archivePath
	task.FailedFiles = failed

	fmt.Printf("success: %s\n", archivePath)
}

func (a *ArchiverService) downloadFile(url, destDir string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed: " + resp.Status)
	}

	filename := filepath.Base(strings.Split(url, "?")[0])
	filename = sanitizeFilename(filename)
	if filename == "" {
		filename = fmt.Sprintf("file_%d", time.Now().UnixNano())
	}

	filePath := filepath.Join(destDir, filename)
	out, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return filePath, err
}

func (a *ArchiverService) createZip(dest string, files []string) error {
	zipFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	for _, file := range files {
		if err := addFileToZip(archive, file); err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = info.Name()
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, file)
	return err
}

func sanitizeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if strings.ContainsRune("\\/:*?\"<>|", r) {
			return -1
		}
		return r
	}, name)
}
