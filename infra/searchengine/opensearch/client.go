package opensearch

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
	"github.com/webitel/webitel-fts/infra/searchengine"
	"io"
	"net/http"
	"time"
)

type OpenSearch struct {
	cli *opensearch.Client
}

func New(hosts []string, username, password string) (*OpenSearch, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: hosts,
		Username:  username,
		Password:  password,
	})

	if err != nil {
		return nil, err
	}

	return &OpenSearch{
		cli: client,
	}, nil
}

func (s *OpenSearch) Shutdown() error {
	return nil
}

func (s *OpenSearch) Test() error {
	req := opensearchapi.PingRequest{
		Pretty: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	res, err := req.Do(ctx, s.cli)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		b, _ := io.ReadAll(res.Body)
		return errors.New(string(b))
	}

	return nil
}

func (s *OpenSearch) Insert(ctx context.Context, id string, index string, body []byte) error {
	document := bytes.NewReader(body)
	req := opensearchapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       document,
	}
	insertResponse, err := req.Do(ctx, s.cli)
	if err != nil {
		return err
	}
	defer insertResponse.Body.Close()

	if insertResponse.IsError() {
		res, _ := io.ReadAll(insertResponse.Body)
		return errors.New(string(res))
	}

	return nil
}

func (s *OpenSearch) Update(ctx context.Context, id string, index string, body []byte) error {

	doc := make(map[string]json.RawMessage)
	doc["doc"] = body

	body, _ = json.Marshal(doc)

	document := bytes.NewReader(body)
	req := opensearchapi.UpdateRequest{
		Index:      index,
		DocumentID: id,
		Body:       document,
	}
	insertResponse, err := req.Do(ctx, s.cli)
	if err != nil {
		return err
	}
	defer insertResponse.Body.Close()

	if insertResponse.IsError() {
		res, _ := io.ReadAll(insertResponse.Body)
		return errors.New(string(res))
	}

	return nil
}

func (s *OpenSearch) Delete(ctx context.Context, id string, index string) error {
	del := opensearchapi.DeleteRequest{
		Index:      index,
		DocumentID: id,
	}
	deleteResponse, err := del.Do(context.Background(), s.cli)
	if err != nil {
		return err
	}
	defer deleteResponse.Body.Close()

	if deleteResponse.IsError() {
		res, _ := io.ReadAll(deleteResponse.Body)
		return errors.New(string(res))
	}

	return nil
}

func (s *OpenSearch) Template(ctx context.Context, name string, body []byte) error {
	document := bytes.NewReader(body)
	req := opensearchapi.IndicesPutIndexTemplateRequest{
		Body: document,
		Name: name,
		//Create: opensearchapi.BoolPtr(true),
	}
	result, err := req.Do(ctx, s.cli)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	if result.IsError() {
		res, _ := io.ReadAll(result.Body)
		return errors.New(string(res))
	}

	return nil
}

type Highlight map[string][]string

type Hits struct {
	Index     string         `json:"_index"`
	Id        string         `json:"_id"`
	Highlight Highlight      `json:"highlight"`
	Source    map[string]any `json:"_source"`
}

type ResponseHits struct {
	Hits []Hits `json:"hits"`
}

type Response struct {
	Hits ResponseHits `json:"hits"`
}

func (s *OpenSearch) Search(ctx context.Context, IndexName []string, text string, limit int) ([]searchengine.SearchResult, error) {

	q := map[string]any{
		"size": limit,
		"sort": []map[string]any{
			{
				"_score": map[string]any{
					"order": "desc",
				},
			},
		},
		"_source": map[string]any{},
		"highlight": map[string]any{
			"fields": map[string]any{
				"*": map[string]any{},
			},
			"require_field_match": false,
			"pre_tags":            []string{"<strong>"},
			"post_tags":           []string{"</strong>"},
		},
		"stored_fields": []string{"*"},
		"query": map[string]any{
			"bool": map[string]any{
				"filter": []map[string]any{
					{
						"query_string": map[string]any{
							"query": text,
						},
					},
				},
			},
		},
	}
	data, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}

	// Search for the document.
	content := bytes.NewReader(data)

	search := opensearchapi.SearchRequest{
		Index: IndexName,
		Body:  content,
	}

	searchResponse, err := search.Do(ctx, s.cli)
	if err != nil {
		return nil, err
	}
	defer searchResponse.Body.Close()

	if searchResponse.IsError() {
		res, _ := io.ReadAll(searchResponse.Body)
		return nil, errors.New(string(res))
	}

	t, _ := io.ReadAll(searchResponse.Body)
	var res Response

	err = json.Unmarshal(t, &res)
	if err != nil {
		return nil, err
	}

	var response []searchengine.SearchResult

	for _, v := range res.Hits.Hits {
		t := ""
		if len(v.Highlight) != 0 {
			for k, v := range v.Highlight {
				t += fmt.Sprintf("%s: %v", k, v[0])
			}
		} else {
			var ok bool
			if _, ok = v.Source["name"]; ok {
				t = fmt.Sprintf("%v", v.Source["name"])
			} else {
				t = "TODO"
			}
		}
		response = append(response, searchengine.SearchResult{
			Index: v.Index,
			Id:    v.Id,
			Text:  t,
		})
	}

	return response, nil
}
