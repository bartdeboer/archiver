package archiver

import (
	"fmt"
	"strings"

	"github.com/bartdeboer/archiver/codecs/codec"
	"github.com/bartdeboer/archiver/codecs/targz"
	"github.com/bartdeboer/archiver/codecs/zip"
)

func NewArchiverByPlatform(platform string) codec.Codec {
	switch platform {
	case "windows":
		return zip.New()
	default:
		return targz.New()
	}
}

func NewArchiverByType(archType string) (codec.Codec, error) {
	switch archType {
	case "zip":
		return zip.New(), nil
	case "targz":
		return targz.New(), nil
	default:
		return nil, fmt.Errorf("unsupported type")
	}
}

func NewArchiverByFilename(file string) (codec.Codec, error) {
	switch {
	case strings.HasSuffix(file, ".zip"):
		return zip.New(), nil
	case strings.HasSuffix(file, ".tar.gz"):
		return targz.New(), nil
	default:
		return nil, fmt.Errorf("unsupported file extension")
	}
}

func AppendExtensionByPlatform(name, platform string) string {
	return NewArchiverByPlatform(platform).AppendExtension(name)
}

func Create(archivePath string, files []codec.ArchiveMap) error {
	arch, err := NewArchiverByFilename(archivePath)
	if err != nil {
		return err
	}
	return arch.Create(archivePath, files)
}

func Extract(archivePath, destDir string) error {
	arch, err := NewArchiverByFilename(archivePath)
	if err != nil {
		return err
	}
	return arch.Extract(archivePath, destDir)
}

func CreateZip(archivePath string, files []codec.ArchiveMap) error {
	return zip.New().Create(archivePath, files)
}

func ExtractZip(archivePath, destDir string) error {
	return zip.New().Extract(archivePath, destDir)
}

func CreateTARGZ(archivePath string, files []codec.ArchiveMap) error {
	return targz.New().Create(archivePath, files)
}

func ExtractTARGZ(archivePath, destDir string) error {
	return targz.New().Extract(archivePath, destDir)
}
