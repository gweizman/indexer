package persistent_storage

import (
	"log"
	"path/filepath"

	"github.com/gocql/gocql"
)

type FileContent struct {
	Data      []byte
	CanBeCode bool
}

type StoredFile struct {
	Path    string
	Name    string
	Content map[string]FileContent
}

type IndexStorage struct {
	dbSession *gocql.Session
	project   string
	basepath  string
}

func NewStorage(project string, basepath string) (*IndexStorage, error) {
	return &IndexStorage{
		nil,
		project,
		basepath,
	}, nil

	cluster := gocql.NewCluster("127.0.0.1")
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	storage := IndexStorage{
		dbSession: session,
		project:   project,
		basepath:  basepath,
	}
	return &storage, nil
}

func (e *IndexStorage) CreateFile(full_path string) (StoredFile, error) {
	dir, file := filepath.Split(full_path)
	reldir, err := filepath.Rel(e.basepath, dir)
	if err != nil {
		return StoredFile{}, err
	}

	return StoredFile{
		Path:    reldir,
		Name:    file,
		Content: make(map[string]FileContent),
	}, nil
}

func (f *StoredFile) AddContentVersion(content []byte, version string) error {
	content_to_print := ""
	if len(content) > 10 {
		content_to_print = string(content[:10])
	} else {
		content_to_print = string(content)
	}
	log.Printf("%s %s %s", f.Path, f.Name, content_to_print)

	f.Content[version] = FileContent{
		Data:      content,
		CanBeCode: true,
	}

	return nil
}

func (f *StoredFile) AddContent(content []byte) error {
	f.AddContentVersion(content, "Original")
	return nil
}

func (f *StoredFile) InsertDefinition(name string, line int, kind string, language string, parent string, parentKind string, pattern string, signature string, fileLimited bool) error {
	log.Printf("%s: %s %d %s %s %s %s %s %s %t", filepath.Join(f.Path, f.Name), name, line, kind, language, parent, parentKind, pattern, signature, fileLimited)
	return nil
}
