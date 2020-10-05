package elasticsearch

import (
	"context"
	"github.com/olivere/elastic"
	"time"
)

var (
	Client esClientInterface = &esClient{}
)

type esClientInterface interface {
	setClient(c *elastic.Client)
	Index(interface{}) (*elastic.IndexResponse, error)
}

type esClient struct {
	client *elastic.Client
}

func Init() {
	client, err := elastic.NewClient(
		elastic.SetURL("http://127.0.01:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		//elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		//elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
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
func (ec *esClient) Index(i interface{}) (*elastic.IndexResponse, error) {
	ctx := context.Background()
	return ec.client.Index().Do(ctx)
}
