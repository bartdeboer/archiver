package archiver

import (
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

func NewArchiverByType(platform string) codec.Codec {
	switch platform {
	case "zip":
		return zip.New()
	case "targz":
		return targz.New()
	}
	return nil
}

func CreateZip(archivePath string, files []codec.ArchiveMap) error {
	return zip.New().Create(archivePath, files)
}

func ExtractZip(filePath, destDir string) error {
	return zip.New().Extract(filePath, destDir)
}

func CreateTARGZ(archivePath string, files []codec.ArchiveMap) error {
	return targz.New().Create(archivePath, files)
}

func ExtractTARGZ(filePath, destDir string) error {
	return targz.New().Extract(filePath, destDir)
}
