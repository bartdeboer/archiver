package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"

	"github.com/bartdeboer/archiver/codecs/codec"
)

// createTarGz creates a .tar.gz archive at archivePath containing the files specified in files.
func CreateTarGz(archivePath string, files []codec.ArchiveMap) error {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

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

		hdr := &tar.Header{
			Name: fileMap.Dest, // Use the target path specified in FilePair
			Mode: int64(fileInfo.Mode()),
			Size: fileInfo.Size(),
		}

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if _, err := io.Copy(tw, srcFile); err != nil {
			return err
		}
	}

	if err := tw.Close(); err != nil {
		return err
	}
	if err := gzw.Close(); err != nil {
		return err
	}

	return os.WriteFile(archivePath, buf.Bytes(), 0644)
}
