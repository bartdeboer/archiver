package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type TarGz struct{}

func New() *TarGz {
	return &TarGz{}
}

const Extension string = ".tar.gz"

func (c *TarGz) AppendExtension(name string) string {
	return name + Extension
}

// create creates a .tar.gz archive at archivePath containing the files specified in files.
func (c *TarGz) Create(archivePath string, files map[string]string) error {
	var buf bytes.Buffer
	gzw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gzw)

	for src, dest := range files {
		srcFile, err := os.Open(src)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		fileInfo, err := srcFile.Stat()
		if err != nil {
			return err
		}

		hdr := &tar.Header{
			Name: dest, // Use the target path specified in FilePair
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

// extract handles the extraction of .tar.gz files.
func (c *TarGz) Extract(archivePath, destDir string) error {

	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("Error creating gzip.NewReader: %v", err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			return fmt.Errorf("Error reading tarReader.Next: %v", err)
		case header == nil:
			continue
		}

		target := filepath.Join(destDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:

			fmt.Printf("Creating directory %s\n", target)

			if err := os.MkdirAll(target, os.FileMode(header.Mode)); err != nil {
				return fmt.Errorf("Error creating directory: %v", err)
			}

		case tar.TypeReg:

			fmt.Printf("Creating file %s\n", target)

			outFile, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("Error creating file: %v", err)
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return fmt.Errorf("Error extracting into file: %v", err)
			}
			outFile.Close()
		}
	}
}
