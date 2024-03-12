package codec

type ArchiveMap struct {
	Src  string
	Dest string
}

// type ArchiveMap interface {
// 	Src() string
// 	Dest() string
// }

type Codec interface {
	Extract(archivePath, destDir string) error
	Create(archivePath string, files []ArchiveMap) error
	AppendExtension(name string) string
}
