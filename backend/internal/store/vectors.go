package store

import (
	"context"
	"encoding/json"
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

type GraphQLResponse struct {
	Data struct {
		Get struct {
			Book []struct {
				Additional struct {
					ID string `json:"id"`
				} `json:"_additional"`
			} `json:"Book"`
		} `json:"Get"`
	} `json:"data"`
}

type IDResponse struct {
	Id string `json:"id"`
}

type SuccessfullyAPIOperation struct {
	Message string `json:"message"`
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
		// if the ok is true than the chapter exits
		if ok {
			log.Printf("chapter: %s already exists in weaviate", doc.Chapter)
			return nil, fmt.Errorf("error: chapter %s %w", doc.Chapter, ErrChapterAlreadyExists)
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

func (d *VectorsStore) GetClosestVectors(ctx context.Context, query string) ([]*Document, error) {
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

func (d *VectorsStore) GetObjectIDByChapter(ctx context.Context, query string) (*IDResponse, error) {

	result, err := d.client.GraphQL().Get().
		WithClassName("Book").
		WithWhere(
			filters.Where().
				WithPath([]string{"chapter"}).
				WithOperator(filters.Equal).
				WithValueText(query),
		).
		WithFields(
			graphql.Field{
				Name: "_additional",
				Fields: []graphql.Field{
					{
						Name: "id",
					},
				},
			},
		).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get object by chapter %s: %w", query, err)
	}

	// Unmarshal the response into the struct
	jsonBytes, err := result.MarshalBinary()
	if err != nil {
		return nil, err
	}

	var response GraphQLResponse
	if err := json.Unmarshal(jsonBytes, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	if len(response.Data.Get.Book) > 0 {
		objIDResponse := &IDResponse{
			Id: response.Data.Get.Book[0].Additional.ID,
		}
		fmt.Println(objIDResponse)

		return objIDResponse, nil
	}

	return nil, fmt.Errorf("no object found for chapter: %s", query)
}

func (d *VectorsStore) DeleteObjectWithID(ctx context.Context, idToDelete string) (*SuccessfullyAPIOperation, error) {
	obj, err := d.client.Data().ObjectsGetter().
		WithClassName("Book").
		WithID(idToDelete).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s retrieving object with id %s", err.Error(), idToDelete)
	}
	if obj == nil {
		return nil, fmt.Errorf("object with id %s does not exist", idToDelete)
	}

	err = d.client.Data().
		Deleter().
		WithClassName("Book").
		WithID(idToDelete).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("error deleting object with id %s: %w", idToDelete, err)
	}

	response := &SuccessfullyAPIOperation{
		Message: "Object deleted successfully ",
	}

	return response, nil
}

func (d *VectorsStore) UpdateObjectWithID(ctx context.Context, updatedDocuments Document, idToUpdate string) (*SuccessfullyAPIOperation, error) {
	obj, err := d.client.Data().ObjectsGetter().
		WithClassName("Book").
		WithID(idToUpdate).
		Do(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s retrieving object with id %s", err.Error(), idToUpdate)
	}
	if obj == nil {
		return nil, fmt.Errorf("error updating object with id %s, it does not exist: %w", idToUpdate, ErrNotFound)
	}

	err = d.client.Data().
		Updater().
		WithClassName("Book").
		WithID(idToUpdate).
		WithProperties(updatedDocuments).
		Do(ctx)

	if err != nil {
		return nil, fmt.Errorf("error updating object with id %s: %w", idToUpdate, err)
	}

	response := &SuccessfullyAPIOperation{
		Message: "Object updated successfully",
	}

	return response, nil
}

// there is no endpoint for this action in the api
func (d *VectorsStore) DeleteChapterWithChapterName(ctx context.Context, chapterName string) (*SuccessfullyAPIOperation, error) {

	ok, err := d.chapterExists(ctx, chapterName)
	if err != nil {
		return nil, fmt.Errorf("error checking if a chapter already exits in weaviate: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("error can not delete chapter, chapter %s does not exits: %w", chapterName, ErrNotFound)
	}
	result, err := d.client.Batch().
		ObjectsBatchDeleter().
		WithClassName("Book").
		WithWhere(
			filters.Where().
				WithPath([]string{"chapter"}).
				WithOperator(filters.Equal).
				WithValueText(chapterName),
		).
		WithOutput("verbose").
		Do(ctx)

	if err != nil {
		return nil, err
	}

	fmt.Print(result)

	response := &SuccessfullyAPIOperation{
		Message: *result.Output,
	}

	return response, nil
}

func parserGraphQLResponseToResponse(res *models.GraphQLResponse) ([]*Document, error) {
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
		return nil, ErrNotFound
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

	documents := make([]*Document, 0, len(chapterMap))
	for _, doc := range chapterMap {
		documents = append(documents, doc)
	}

	return documents, nil
}

// false = chapter not found / true = chapter found
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

		return false, err
	}
	getData, ok := response.Data["Get"].(map[string]any)
	if !ok {
		return false, fmt.Errorf("invalid response structure: missing 'Get'")
	}
	rawBooks, ok := getData["Book"].([]any)
	if !ok || len(rawBooks) == 0 {
		return false, nil
	}
	return true, nil
}
