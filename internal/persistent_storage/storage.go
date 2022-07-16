package persistent_storage

import (
	"log"
	"path/filepath"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type FileContent struct {
	Data      []byte
	CanBeCode bool
}

type IndexStorage struct {
	db       Db
	project  string
	basepath string
}

type StoredFile struct {
	Idx     *IndexStorage
	Path    string
	Name    string
	Content map[string]FileContent
}

func NewStorage(project string, basepath string, db Db) (*IndexStorage, error) {
	session := db.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.Run("MATCH (n) DETACH DELETE n", map[string]interface{}{}) // TODO: Delete me!
	if err != nil {
		return nil, err
	}

	_, err = session.Run("MERGE (a:Project{name: $project})", map[string]interface{}{"project": project})
	if err != nil {
		return nil, err
	}

	_, err = session.Run(""+
		"MATCH (p:Project) "+
		"WHERE p.name = $project "+
		"MERGE (f:Dir {path: $path})-[:BELONGS_TO]->(p)",
		map[string]interface{}{"project": project, "path": "."},
	)
	if err != nil {
		return nil, err
	}

	storage := IndexStorage{
		db:       db,
		project:  project,
		basepath: basepath,
	}
	return &storage, nil
}

func (e *IndexStorage) CreateDir(full_path string) error {
	reldir, err := filepath.Rel(e.basepath, full_path)
	if err != nil {
		return err
	}
	session := e.db.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.Run(""+
		"MATCH (m:Dir) "+
		"WHERE m.path = $parent "+
		"MERGE (m)-[:CONTAINS]->(f:Dir {path: $path, project: $project})",
		map[string]interface{}{"project": e.project, "path": reldir, "parent": filepath.Join(reldir, "../")},
	)
	if err != nil {
		return err
	}

	return nil
}

func (e *IndexStorage) CreateFile(full_path string) (StoredFile, error) {
	dir, file := filepath.Split(full_path)
	reldir, err := filepath.Rel(e.basepath, dir)
	if err != nil {
		return StoredFile{}, err
	}
	session := e.db.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err = session.Run(""+
		"MATCH (m:Dir) "+
		"WHERE m.path = $path "+
		"CREATE (m)-[:CONTAINS]->(f:File {path: $path, name: $name, project: $project})",
		map[string]interface{}{"project": e.project, "path": reldir, "name": file},
	)
	if err != nil {
		log.Printf("%s %s %s %s", err, e.project, reldir, file)
		return StoredFile{}, err
	}

	return StoredFile{
		Idx:     e,
		Path:    reldir,
		Name:    file,
		Content: make(map[string]FileContent),
	}, nil
}

func (f *StoredFile) AddContentVersion(content []byte, version string) error {
	session := f.Idx.db.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

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

	_, err := session.Run(""+
		"MATCH (f:File) "+
		"WHERE f.path = $file_path AND f.name = $file_name AND f.project = $project "+
		"CREATE (f)-[:CONTENT]->(c:Content {project: $project, version: $version, data: $data})",
		map[string]interface{}{"file_name": f.Name, "file_path": f.Path, "project": f.Idx.project, "version": version, "data": string(content)})
	if err != nil {
		return err
	}

	return nil
}

func (f *StoredFile) AddContent(content []byte) error {
	f.AddContentVersion(content, "original")
	return nil
}

func (f *StoredFile) InsertDefinition(name string, line int, kind string, language string, parent string, parentKind string, pattern string, signature string, fileLimited bool) error {
	session := f.Idx.db.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	_, err := session.Run(""+
		"MATCH (f:File) "+
		"WHERE f.path = $file_path AND f.name = $file_name AND f.project = $project "+
		"CREATE (f)-[:DEFINES]->(d:Definition {project: $project, name: $name, line: $line, kind: $kind, language: $language, parent: $parent, parentKind: $parentKind, pattern: $pattern, signature: $signature, fileLimited: $fileLimited})",
		map[string]interface{}{"file_name": f.Name, "file_path": f.Path, "project": f.Idx.project, "name": name, "line": line, "kind": kind, "language": language, "parent": parent, "parentKind": parentKind, "pattern": pattern, "signature": signature, "fileLimited": fileLimited})
	if err != nil {
		return err
	}

	log.Printf("%s: %s %d %s %s %s %s %s %s %t", filepath.Join(f.Path, f.Name), name, line, kind, language, parent, parentKind, pattern, signature, fileLimited)
	return nil
}
