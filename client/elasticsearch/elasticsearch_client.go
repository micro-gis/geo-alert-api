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
}

type esClient struct {
	client *elastic.Client
}

func Init() {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.01:9200"),
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
