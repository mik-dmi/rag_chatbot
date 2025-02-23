package store

import (
	"context"
	"fmt"
	"log"
	"strings"

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

type VectorCreatedResponse struct {
	ChaptersCreated []string `json:"list_of_chapters_created"`
}

func (d *VectorsStore) CreateVectors(ctx context.Context, data *RagData) (*VectorCreatedResponse, error) {
	var objects []*models.Object
	var chaptersCreated []string

	// Check if Chapter of data already exists

	for _, doc := range data.Documents {
		ok, err := d.chapterExists(ctx, doc.Chapter)
		if err != nil {
			return nil, fmt.Errorf("error checking if a chapter already exits in weaviate: %w", err)
		}
		// if the ok  is false than the chapter
		if ok {
			log.Printf("chapter: %s already exists in weaviate", doc.Chapter)
			return nil, fmt.Errorf("error chapter %s already exits in weaviate", doc.Chapter)
		}

		var subs []Subsection
		for _, subsection := range doc.Subsections {
			subs = append(subs, Subsection{
				Title:   subsection.Title,
				Content: subsection.Content,
			})
		}

		obj := &models.Object{
			Class: "Book",
			Properties: map[string]interface{}{
				"chapter":     doc.Chapter,
				"subsections": subs,
			},
		}
		objects = append(objects, obj)
		chaptersCreated = append(chaptersCreated, doc.Chapter)

	}
	_, err := d.client.Batch().ObjectsBatcher().WithObjects(objects...).Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("batch insert failed: %w", err)
	}

	jsonChapters := VectorCreatedResponse{
		ChaptersCreated: chaptersCreated,
	}

	return &jsonChapters, nil
}

func (d *VectorsStore) GetClosestVectors(ctx context.Context, query string) ([]Document, error) {
	maxDistance := float32(0.5) //max similarity threshold
	graphQLResponse, err := d.client.GraphQL().Get().
		WithClassName("Book").
		WithFields(
			graphql.Field{Name: "chapter"},
			graphql.Field{
				Name: "subsections",
				Fields: []graphql.Field{
					{Name: "title"},
					{Name: "content"},
				},
			},
		).
		WithNearText(d.client.GraphQL().NearTextArgBuilder().
			WithConcepts([]string{query}).
			WithDistance(maxDistance)).
		WithLimit(5).
		Do(ctx)

	if err != nil {
		return nil, err
	}

	//
	jsonBytes, err := graphQLResponse.MarshalBinary()
	if err != nil {
		return nil, err
	}

	jsonStr := string(jsonBytes)
	fmt.Println("Response from jsonStr:", jsonStr)

	//

	response, err := parserGraphQLResponseToResponse(graphQLResponse)
	if err != nil {
		return nil, err
	}
	fmt.Println("Response from parserGraphQLResponseToResponse:", response)

	return response, nil
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

func parserGraphQLResponseToResponse(res *models.GraphQLResponse) ([]Document, error) {
	if len(res.Errors) > 0 {
		messages := make([]string, 0, len(res.Errors))
		for _, e := range res.Errors {
			messages = append(messages, e.Message)
		}
		return nil, fmt.Errorf("GraphQL response error: %s", strings.Join(messages, ", "))
	}

	getData, ok := res.Data["Get"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid response structure: missing 'Get' key")
	}

	rawBooks, ok := getData["Book"].([]any)
	if !ok || len(rawBooks) == 0 {
		return nil, fmt.Errorf("no Book data found")
	}

	// Group results by chapter.
	chapterMap := make(map[string]*Document)
	for _, item := range rawBooks {
		itemMap, ok := item.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("invalid item format in response")
		}

		/*
			Check if the fields are nested under "properties". If yes, use that.
			var props map[string]any
			if p, exists := itemMap["properties"].(map[string]any); exists {
				fmt.Println("Has properties")
				props = p
			} else {
				fmt.Println("Does not have properties")
				props = itemMap
			}*/

		chapter, _ := itemMap["chapter"].(string)

		subsectionsRaw, ok := itemMap["subsections"].([]any)
		if !ok {
			subsectionsRaw = []any{}
		}
		var subs []Subsection
		for _, subItem := range subsectionsRaw {
			subMap, ok := subItem.(map[string]any)
			if !ok {
				continue
			}
			title, _ := subMap["title"].(string)
			content, _ := subMap["content"].(string)
			subs = append(subs, Subsection{
				Title:   title,
				Content: content,
			})
		}

		if doc, exists := chapterMap[chapter]; exists {
			doc.Subsections = append(doc.Subsections, subs...)
		} else {
			chapterMap[chapter] = &Document{
				Chapter:     chapter,
				Subsections: subs,
			}
		}
	}

	documents := make([]Document, 0, len(chapterMap))
	for _, doc := range chapterMap {
		documents = append(documents, *doc)
	}

	return documents, nil
}

func (d *VectorsStore) chapterExists(ctx context.Context, chapter string) (bool, error) {
	response, err := d.client.GraphQL().Get().
		WithClassName("Book").
		WithFields(graphql.Field{Name: "chapter"}).
		WithWhere(filters.Where().
			WithPath([]string{"chapter"}).
			WithOperator(filters.Equal).
			WithValueString(chapter)).
		WithLimit(1).
		Do(ctx)
	if err != nil {
		fmt.Println("Yes 1")
		return false, err
	}

	getData, ok := response.Data["Get"].(map[string]any)
	if !ok {
		fmt.Println("Yes 2")
		return false, fmt.Errorf("invalid response structure: missing 'Get'")
	}

	rawBooks, ok := getData["Book"].([]any)
	//  no chapter found
	if !ok || len(rawBooks) == 0 {
		return false, nil
	}
	//  chapter found
	return true, nil
}
