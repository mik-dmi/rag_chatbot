package store

import (
	"context"

	"github.com/tmc/langchaingo/schema"
	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate/entities/models"
)

type RagData struct {
	UserID    string     `json:"userid"`
	Documents []Document `json:"document"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}
type Document struct {
	Chapter    string       `json:"chapter"`
	Subsection []Subsection `json:"subsection"`
}
type Subsection struct {
	Title   []string `json:"title"`
	Content []string `json:"content"`
}
type VectorsStore struct {
	client *weaviate.Client
}

func (d *VectorsStore) CreateVectors(ctx context.Context, data *RagData) error {

	var docs []schema.Document
	for i, doc := range data.Documents {
		docs = append(docs, schema.Document{
			PageContent: doc.Content,
			Metadata: map[string]interface{}{
				"source": doc.Title,
			},
		})

		// Define the collection
		classObj := &models.Class{
			Class:      "Question",
			Vectorizer: "text2vec-transformers",
		}

		// add the collection
		err = client.Schema().ClassCreator().WithClass(classObj).Do(context.Background())
		if err != nil {
			panic(err)
		}

		// convert items into a slice of models.Object
		objects := make([]*models.Object, len(docs))

		objects[i] = &models.Object{
			Class: "Question",
			Properties: map[string]any{
				"chapter": doc.Chapter,
				"title":   doc.Title,
				"content": doc.Content,
			},
		}

		// batch write items
		batchRes, err := client.Batch().ObjectsBatcher().WithObjects(objects...).Do(context.Background())
		if err != nil {
			panic(err)
		}

	}

	return nil
}

func (d *VectorsStore) DeleteVectors(ctx context.Context, deleteDocument string) error {
	return nil
}
