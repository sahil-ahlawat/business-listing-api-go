// File: fitness/test/elasticsearch/elastic_test.go
package elastic_test

import (
	"context"
	"testing"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

func TestElasticSearchConnection(t *testing.T) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		t.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	res, err := es.Info()
	if err != nil || res.IsError() {
		t.Fatalf("Elasticsearch connection failed: %v", err)
	}
}