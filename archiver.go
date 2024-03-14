package archiver

import (
	"fmt"
	"strings"

	"github.com/bartdeboer/archiver/targz"
	"github.com/bartdeboer/archiver/zip"
)

type Archiver interface {
	Extract(archivePath, destDir string) error
	Create(archivePath string, files map[string]string) error
	AppendExtension(name string) string
}

type implementation struct {
	create    func() Archiver
	extension string
}

var impls map[string]implementation

func RegisterArchiver(name string, create func() Archiver, extension string) {
	impls[name] = implementation{
		create,
		extension,
	}
}

func init() {
	impls = make(map[string]implementation)
	RegisterArchiver("windows", func() Archiver { return zip.New() }, zip.Extension)
	RegisterArchiver("targz", func() Archiver { return targz.New() }, targz.Extension)
}

func New(archType string) (Archiver, error) {
	if impl, exists := impls[archType]; exists {
		return impl.create(), nil
	}
	return nil, fmt.Errorf("unsupported type")
}

func NewByFilename(file string) (Archiver, error) {
	for _, impl := range impls {
		if strings.HasSuffix(file, impl.extension) {
			return impl.create(), nil
		}
	}
	return nil, fmt.Errorf("unsupported file extension")
}

func NewArchiverByPlatform(platform string) Archiver {
	switch platform {
	case "windows":
		return zip.New()
	default:
		return targz.New()
	}
}

func AppendExtensionByPlatform(name, platform string) string {
	return NewArchiverByPlatform(platform).AppendExtension(name)
}

func Create(archivePath string, files map[string]string) error {
	arch, err := NewByFilename(archivePath)
	if err != nil {
		return err
	}
	return arch.Create(archivePath, files)
}

func Extract(archivePath, destDir string) error {
	arch, err := NewByFilename(archivePath)
	if err != nil {
		return err
	}
	return arch.Extract(archivePath, destDir)
}

func CreateZip(archivePath string, files map[string]string) error {
	return zip.New().Create(archivePath, files)
}

func ExtractZip(archivePath, destDir string) error {
	return zip.New().Extract(archivePath, destDir)
}

func CreateTarGz(archivePath string, files map[string]string) error {
	return targz.New().Create(archivePath, files)
}

func ExtractTarGz(archivePath, destDir string) error {
	return targz.New().Extract(archivePath, destDir)
}
