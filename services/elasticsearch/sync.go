package elasticsearch

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"fitness/models"
	"fitness/services/nats"

	"github.com/nats-io/nats.go"
)

func Sync() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Syncing with Elasticsearch...")
		nats.PollMessages("elastic-search-sync", "es-sync-worker", func(msg *nats.Msg) bool {
			var listing models.Listing
			if err := json.Unmarshal(msg.Data, &listing); err != nil {
				log.Printf("Error unmarshalling message: %v", err)
				return true // Acknowledge the message to avoid redelivery
			}

			if err := indexListing(listing); err != nil {
				log.Printf("Error indexing listing: %v", err)
				return false // Do not acknowledge the message, so it can be redelivered
			}

			return true // Acknowledge the message
		}, 100)
	}
}

func indexListing(listing models.Listing) error {
	ctx := context.Background()
	listingJSON, err := json.Marshal(listing)
	if err != nil {
		return err
	}

	res, err := esClient.Index(
		"listings",
		strings.NewReader(string(listingJSON)),
		esClient.Index.WithDocumentID(listing.Slug),
		esClient.Index.WithContext(ctx),
		esClient.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("error indexing document ID=%s: %s", listing.Slug, res.String())
	}

	return nil
}
