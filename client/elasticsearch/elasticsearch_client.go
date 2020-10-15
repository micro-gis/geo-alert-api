package elasticsearch

import (
	"context"
	"fmt"
	"github.com/micro-gis/utils/logger"
	"github.com/olivere/elastic"
	"time"
)

var (
	Client esClientInterface = &esClient{}
)

type esClientInterface interface {
	setClient(c *elastic.Client)
	Index(string, string, interface{}) (*elastic.IndexResponse, error)
	Get(string, string, string) (*elastic.GetResult, error)
	Delete(string, string, string) (*elastic.DeleteResponse, error)
	Search(string, elastic.Query) (*elastic.SearchResult, error)
	Upsert(string, string, interface{}, string) (*elastic.UpdateResponse, error)
}

type esClient struct {
	client *elastic.Client
}

func Init() {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.01:9100"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetErrorLog(logger.GetLogger()),
		elastic.SetInfoLog(logger.GetLogger()),
		//elastic.SetHeaders(http.Header{
		//	"X-Caller-Id": []string{"..."},
		//}),
	)

	if err != nil {
		panic(err)
	}
	Client.setClient(client)
}
func (ec *esClient) setClient(c *elastic.Client) {
	ec.client = c
}
func (ec *esClient) Index(index string, doctype string, doc interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	result, err := ec.client.Index().
		Type(doctype).
		Index(index).
		BodyJson(doc).
		Do(ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to index document in index : %s", index), err)
		return nil, err
	}
	return result, nil
}

func (ec *esClient) Get(index string, doctype string, id string) (*elastic.GetResult, error) {
	ctx := context.Background()
	result, err := ec.client.Get().Index(index).Type(doctype).Id(id).Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to get id %s", id), err)
		return nil, err
	}
	return result, nil
}

func (ec *esClient) Search(index string, query elastic.Query) (*elastic.SearchResult, error) {
	ctx := context.Background()
	result, err := ec.client.Search(index).
		Query(query).
		RestTotalHitsAsInt(true).
		Do(ctx)
	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to search documents in index %s", index), err)
		return nil, err
	}
	return result, nil
}

func (c *esClient) Delete(index string, docType string, id string) (*elastic.DeleteResponse, error) {
	ctx := context.Background()
	result, err := c.client.Delete().
		Index(index).
		Type(docType).
		Id(id).
		Do(ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to delete document id %s in index %s", id, index), err)
		return nil, err
	}
	return result, nil
}

func (c *esClient) Upsert(index string, docType string, doc interface{}, id string) (*elastic.UpdateResponse, error) {
	ctx := context.Background()
	update, err := c.client.Update().
		Index(index).
		Type(docType).
		Id(id).
		Doc(doc).
		DocAsUpsert(true).
		Do(ctx)

	if err != nil {
		logger.Error(fmt.Sprintf("error when trying to update document of id  %s in index %s", id, index), err)
		return nil, err
	}
	return update, nil
}
