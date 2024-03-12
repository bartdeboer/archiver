package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"

	"github.com/bartdeboer/archiver/codecs/codec"
)

type Codec struct{}

func New() *Codec {
	return &Codec{}
}

func (c *Codec) AppendExtension(name string) string {
	return name + ".zip"
}

// create creates a zip archive at archivePath containing the files specified in archiveMaps.
func (c *Codec) Create(archivePath string, files []codec.ArchiveMap) error {
	zipFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	for _, fileMap := range files {
		srcFile, err := os.Open(fileMap.Src)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		fileInfo, err := srcFile.Stat()
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(fileInfo)
		if err != nil {
			return err
		}
		header.Name = fileMap.Dest  // Use the destination path specified in ArchiveMap
		header.Method = zip.Deflate // Use compression

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if _, err := io.Copy(writer, srcFile); err != nil {
			return err
		}
	}

	return nil
}

// extract handles the extraction of .zip files.
func (c *Codec) Extract(archivePath, destDir string) error {
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		fPath := filepath.Join(destDir, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(fPath, os.ModePerm)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(fPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.OpenFile(fPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		rc, err := file.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}
	return nil
}
