package zip

import (
	"archive/zip"
	"io"
	"os"

	"github.com/bartdeboer/archiver/codecs/codec"
)

type Codec struct{}

func New() *Codec {
	return &Codec{}
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
