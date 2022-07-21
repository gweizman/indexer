package persistent_storage

import (
	"context"
	"log"
	"path/filepath"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/pkg/errors"
	"github.com/stevenferrer/solr-go"
)

/*
	General DB wrappers
*/

type Db struct {
	cancel      context.CancelFunc
	ctx         context.Context
	solrClient  solr.Client
	neo4jDriver neo4j.Driver
}

func CreateDb() (*Db, error) {
	ctx, cancel := context.WithCancel(context.Background())
	solrClient := solr.NewJSONClient("http://localhost:8983")

	neo4jDriver, err := neo4j.NewDriver("bolt://localhost:7687", neo4j.BasicAuth("neo4j", "password", ""))
	if err != nil {
		cancel()
		return nil, err
	}

	return &Db{
		cancel:      cancel,
		ctx:         ctx,
		solrClient:  solrClient,
		neo4jDriver: neo4jDriver,
	}, nil
}

func (e *Db) Close() {
	log.Print("Terminating")
	e.neo4jDriver.Close()
	e.cancel()
}

func initSchema(ctx context.Context, collection string, client solr.Client) error {
	fieldTypes := []solr.FieldType{
		// approach #2
		// Refer to https://blog.griddynamics.com/implement-autocomplete-search-for-large-e-commerce-catalogs/
		{
			Name:                 "text_suggest",
			Class:                "solr.TextField",
			PositionIncrementGap: "100",
			IndexAnalyzer: &solr.Analyzer{
				Tokenizer: &solr.Tokenizer{
					Class: "solr.WhitespaceTokenizerFactory",
				},
				Filters: []solr.Filter{
					{
						Class: "solr.LowerCaseFilterFactory",
					},
					{
						Class: "solr.ASCIIFoldingFilterFactory",
					},
					{
						Class:       "solr.EdgeNGramFilterFactory",
						MinGramSize: 1,
						MaxGramSize: 20,
					},
				},
			},
			QueryAnalyzer: &solr.Analyzer{
				Tokenizer: &solr.Tokenizer{
					Class: "solr.WhitespaceTokenizerFactory",
				},
				Filters: []solr.Filter{
					{
						Class: "solr.LowerCaseFilterFactory",
					},
					{
						Class: "solr.ASCIIFoldingFilterFactory",
					},
				},
			},
		},
	}
	err := client.AddFieldTypes(ctx, collection, fieldTypes...)
	if err != nil {
		return err
	}

	// define the fields
	fields := []solr.Field{
		{
			Name: "path",
			Type: "string",
		},
		{
			Name: "project",
			Type: "string",
		},
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "content",
			Type: "text_general",
		},
	}

	err = client.AddFields(ctx, collection, fields...)
	if err != nil {
		return errors.Wrap(err, "add fields")
	}

	return nil
}

const solrCollection = "collection"

func (e *Db) initCore(createCollection bool, shouldInitSchema bool) error {
	if createCollection {
		err := e.solrClient.CreateCollection(e.ctx, solr.NewCollectionParams().Name(solrCollection).NumShards(1).ReplicationFactor(1))
		if err != nil {
			return err
		}
	}

	if shouldInitSchema {
		err := initSchema(e.ctx, solrCollection, e.solrClient)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Db) initConstraints() error {
	session := e.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	var err error = nil

	_, err = session.Run("MATCH (n) DETACH DELETE n", map[string]interface{}{})
	if err != nil {
		return err
	}

	_, err = session.Run("CREATE CONSTRAINT project_name_constraint IF NOT EXISTS FOR (n:Project) REQUIRE n.name IS UNIQUE", map[string]interface{}{})
	if err != nil {
		return err
	}

	_, err = session.Run("CREATE CONSTRAINT dir_path_constraint IF NOT EXISTS FOR (n:Dir) REQUIRE (n.project, n.path) IS UNIQUE", map[string]interface{}{})
	if err != nil {
		return err
	}

	_, err = session.Run("CREATE INDEX file_name_idx IF NOT EXISTS FOR (n:File) ON (n.project, n.name) ", map[string]interface{}{})
	if err != nil {
		return err
	}

	_, err = session.Run("CREATE INDEX def_project_idx FOR (n:Definition) ON (n.project, n.name)", map[string]interface{}{})
	if err != nil {
		return err
	}

	_, err = session.Run("CREATE INDEX def_idx FOR (n:Definition) ON (n.name)", map[string]interface{}{})
	if err != nil {
		return err
	}

	_, err = session.Run("CREATE FULLTEXT INDEX content_project_idx FOR (n:Content) ON EACH [n.data]", map[string]interface{}{}) // TODO: Figure out how to include project in this index
	if err != nil {
		return err
	}

	return err
}

func (e *Db) InitSchemas() error {
	if false { // Not actualy using solr for now
		err := e.initCore(true, true)
		if err != nil {
			return err
		}
	}

	err := e.initConstraints()
	if err != nil {
		return err
	}

	return nil
}

func (e *Db) GetFileContent(project string, path string) (string, bool, error) {
	session := e.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	dir, file := filepath.Split(path)
	dir, err := filepath.Rel("./", dir)
	if err != nil {
		log.Panic(err)
	}

	log.Print(dir)
	log.Print(file)
	output, err := session.Run("MATCH (f:File)-[:CONTENT]->(c:Content) WHERE f.project = $project AND f.path = $path AND f.name = $name RETURN c.data AS data",
		map[string]interface{}{"project": project, "path": filepath.Dir(path), "name": file, "dir": dir}) // TODO: Figure out how to include project in this index
	if err != nil {
		return "", false, err
	}
	if output.Next() {
		content, found := output.Record().Get("data")
		if found {
			return content.(string), found, nil
		}
	}

	return "", false, nil
}

type FileContent struct {
	Project  string  `json:"project"`
	Version  string  `json:"version"`
	FileName string  `json:"file_name"`
	FilePath string  `json:"file_path"`
	Data     string  `json:"data"`
	Score    float64 `json:"score"`
}

// TODO: Don't ignore path_limit
// TODO: Maybe should be the same function as GetFileContent
func (e *Db) SearchFileContent(project string, path_limit string, query string) ([]FileContent, bool, error) {
	session := e.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	output, err := session.Run(""+
		"CALL db.index.fulltext.queryNodes(\"content_project_idx\", $query) YIELD node, score "+
		"MATCH (f:File) "+
		"WHERE (f:File)-[:CONTENT]->(node) AND node.project = $project AND f.project = $project "+
		"RETURN node.data as data, node.project as project, node.version as version, f.name as fileName, f.path as filePath, score as score",
		map[string]interface{}{"project": project, "query": query})
	if err != nil {
		return []FileContent{}, false, err
	}

	var results []FileContent
	for output.Next() {
		record := output.Record()

		project, _ := record.Get("project")
		version, _ := record.Get("version")
		fileName, _ := record.Get("fileName")
		filePath, _ := record.Get("filePath")
		data, _ := record.Get("data")
		score, _ := record.Get("score")

		results = append(results, FileContent{
			Project:  project.(string),
			Version:  version.(string),
			FileName: fileName.(string),
			FilePath: filePath.(string),
			Data:     data.(string),
			Score:    score.(float64),
		})
	}

	return results, true, nil
}

type Definition struct {
	Project     string `json:"project"`
	Name        string `json:"name"`
	Language    string `json:"language"`
	Pattern     string `json:"pattern"`
	Signature   string `json:"signature"`
	FileLimited bool   `json:"file_limited"`
	Parent      string `json:"parent"`
	ParentKind  string `json:"parent_kind"`
	FileName    string `json:"file_name"`
	FilePath    string `json:"file_path"`
	Line        uint   `json:"line"`
}

func (e *Db) GetDefinition(project string, path_limit string, name string) ([]Definition, bool, error) { // TODO: Don't ignore path_limit
	session := e.neo4jDriver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	output, err := session.Run("MATCH (n:Definition)<-[:DEFINES]-(f:File) WHERE n.project = $project AND n.name = $name RETURN n.project as project, n.name as name, "+
		"n.language as language, n.pattern as pattern, n.signature as signature, n.fileLimited as fileLimited, n.parent as parent, n.parentKind as parentKind, n.line as line, f.name as fileName, f.path as filePath",
		map[string]interface{}{"project": project, "name": name})
	if err != nil {
		return []Definition{}, false, err
	}

	var results []Definition
	for output.Next() {
		record := output.Record()

		project, _ := record.Get("project")
		name, _ := record.Get("name")
		language, _ := record.Get("language")
		pattern, _ := record.Get("pattern")
		signature, _ := record.Get("signature")
		fileLimited, _ := record.Get("fileLimited")
		parent, _ := record.Get("parent")
		parentKind, _ := record.Get("parentKind")
		fileName, _ := record.Get("fileName")
		filePath, _ := record.Get("filePath")
		line, _ := record.Get("line")

		results = append(results, Definition{
			Project:     project.(string),
			Name:        name.(string),
			Language:    language.(string),
			Pattern:     pattern.(string),
			Signature:   signature.(string),
			FileLimited: fileLimited.(bool),
			Parent:      parent.(string),
			ParentKind:  parentKind.(string),
			FileName:    fileName.(string),
			FilePath:    filePath.(string),
			Line:        uint(line.(int64)), // TODO: Fix, should be unsigned @ db
		})
	}

	return results, true, nil
}
