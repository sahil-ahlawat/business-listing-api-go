package elasticsearch

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

var (
	esClient *elasticsearch.Client
	ctx      = context.Background()
)

func Init() {
	var err error
	esClient, err = elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
}

func Search(index, query string) ([]byte, error) {
	res, err := esClient.Search(
		esClient.Search.WithContext(ctx),
		esClient.Search.WithIndex(index),
		esClient.Search.WithBody(strings.NewReader(query)),
		esClient.Search.WithTrackTotalHits(true),
		esClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, err
		} else {
			// Print the business logic error details.
			return nil, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
		}
	}

	return ioutil.ReadAll(res.Body)
}