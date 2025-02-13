package store

import (
	"context"
	"fmt"

	"github.com/weaviate/weaviate-go-client/v4/weaviate"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/filters"
	"github.com/weaviate/weaviate-go-client/v4/weaviate/graphql"
	"github.com/weaviate/weaviate/entities/models"
)

type RagData struct {
	UserID    string     `json:"userid"`
	Documents []Document `json:"document"`
	CreatedAt string     `json:"created_at"`
	UpdatedAt string     `json:"updated_at"`
}
type Document struct {
	Chapter     string       `json:"chapter"`
	Subsections []Subsection `json:"subsections"`
}
type Subsection struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
type VectorsStore struct {
	client *weaviate.Client
}

type Query struct {
	UserID      string `json:"user_id"`
	UserMessage string `json:"user_message"`
	SendAt      string `json:"send_at"`
}

func (d *VectorsStore) CreateVectors(ctx context.Context, data *RagData) error {
	classObj := &models.Class{
		Class:      "Book",
		Vectorizer: "text2vec-transformers",
	}

	if err := d.client.Schema().ClassCreator().WithClass(classObj).Do(ctx); err != nil {
		// Depending on your needs, you might want to check if the class already exists
		// and ignore the error if so.
		return fmt.Errorf("failed to create class: %w", err)
	}

	var objects []*models.Object
	for _, doc := range data.Documents {
		for _, subsection := range doc.Subsections {
			obj := &models.Object{
				Class: "Book",
				Properties: map[string]string{
					"chapter": doc.Chapter,
					"title":   subsection.Title,
					"content": subsection.Content,
				},
			}
			objects = append(objects, obj)
		}
	}

	_, err := d.client.Batch().ObjectsBatcher().WithObjects(objects...).Do(ctx)
	if err != nil {
		return fmt.Errorf("batch insert failed: %w", err)
	}

	return nil
}

func (d *VectorsStore) GetClosestVectors(ctx context.Context, query string) (string, error) {
	maxDistance := float32(0.18) //max similarity threshold
	response, err := d.client.GraphQL().Get().
		WithClassName("Book").
		WithFields(

			graphql.Field{Name: "content"},
		).
		WithNearText(d.client.GraphQL().NearTextArgBuilder().
			WithConcepts([]string{query}).
			WithDistance(maxDistance)).
		WithLimit(5).
		Do(ctx)

	if err != nil {
		return "", err
	}
	fmt.Printf("I'm here:\n")
	fmt.Printf("%v", response)

	return "", nil
}

func (d *VectorsStore) DeleteChapterVectors(ctx context.Context, chapterName string) error {

	result, err := d.client.GraphQL().Get().
		WithClassName("Book").
		WithFields(
			graphql.Field{Name: "chapter"},
			graphql.Field{Name: "title"},
			graphql.Field{Name: "content"},
			graphql.Field{
				Name: "_additional",
				Fields: []graphql.Field{
					{Name: "id"},
				},
			}).
		WithWhere(filters.Where().
			WithPath([]string{"chapter"}).
			WithOperator(filters.Equal).
			WithValueText(chapterName)).
		Do(context.Background())
	if err != nil {
		return err
	}

	fmt.Print(result)

	return nil
}
