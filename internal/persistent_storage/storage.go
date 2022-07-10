package persistent_storage

import (
	"github.com/gocql/gocql"
)

type ParsedFile struct {
	Path   string
	Plain  []byte
	Parsed []byte
}

type IndexStorage struct {
	dbSession *gocql.Session
	project   string
}

func NewStorage(project string) (*IndexStorage, error) {
	return nil, nil

	cluster := gocql.NewCluster("127.0.0.1")
	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	storage := IndexStorage{
		dbSession: session,
		project:   project,
	}
	return &storage, nil
}

func (e *IndexStorage) InsertFileVersion(path string, name string, content []byte, version string) {
	return
}

func (e *IndexStorage) InsertFile(path string, name string, content []byte) {
	e.InsertFileVersion(path, name, content, "Original")
}

func (e *IndexStorage) InsertDefinition(file ParsedFile, name string, line int, kind string, language string, parent string, parentKind string, pattern string, signature string, fileLimited bool) {
}
