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
