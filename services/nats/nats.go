// File: fitness/nats/nats.go
package nats

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"fitness/config"
	"fitness/models"

	"github.com/nats-io/nats.go"
)

var Conn *nats.Conn
type Msg = *nats.Msg
var JS nats.JetStreamContext

func Init() {
	config.LoadEnv()

	natsHost := config.GetEnv("NATS_HOST", "localhost")
	natsPort := config.GetEnv("NATS_PORT", "4222")

	var err error
	Conn, err = nats.Connect(natsHost + ":" + natsPort)
	if err != nil {
		log.Fatal("NATS connection failed:", err)
	}

	JS, err = Conn.JetStream()
	if err != nil {
		log.Fatal("JetStream context failed:", err)
	}
}

// CreateStream sets up a stream if it doesn't exist
func CreateStream(streamName, subject string) error {
	_, err := JS.AddStream(&nats.StreamConfig{
		Name:     streamName,
		Subjects: []string{subject},
		Storage:  nats.FileStorage,
	})
	if err != nil && err != nats.ErrStreamNameAlreadyInUse {
		return fmt.Errorf("stream creation failed: %w", err)
	}
	return nil
}

// ListStreams lists all stream names
func ListStreams() ([]string, error) {
	streams := []string{}
	for info := range JS.StreamNames() {
		streams = append(streams, info)
	}
	return streams, nil
}

// ListSubjects for a stream
func ListSubjects(streamName string) ([]string, error) {
	info, err := JS.StreamInfo(streamName)
	if err != nil {
		return nil, err
	}
	return info.Config.Subjects, nil
}

// Publish sends a message to a given subject
func Publish(subject, message string) error {
	_, err := JS.Publish(subject, []byte(message))
	return err
}

// Subscribe sets up a consumer that acknowledges messages manually
func Subscribe(subject string, handler func(msg *nats.Msg)) error {
	_, err := JS.Subscribe(subject, func(m *nats.Msg) {
		handler(m)
	}, nats.ManualAck())
	return err
}

// PollMessages polls a subject for messages every minute and optionally acks them
func PollMessages(subject, durable string, process func(msg *nats.Msg) bool, maxMessages int) error {
	sub, err := JS.PullSubscribe(subject, durable, nats.BindStream(""))
	if err != nil {
		return fmt.Errorf("pull subscription failed: %w", err)
	}

	go func() {
		for {
			msgs, err := sub.Fetch(maxMessages, nats.MaxWait(5*time.Second))
			if err != nil {
				log.Println("No messages fetched or error:", err)
			} else {
				for _, msg := range msgs {
					if process(msg) {
						msg.Ack()
					} else {
						log.Println("Message skipped, not acknowledged")
					}
				}
			}
			time.Sleep(1 * time.Minute)
		}
	}()

	return nil
}

// ProcessListingInsert inserts the listing into MySQL if valid
func ProcessListingInsert(db *sql.DB, msg *nats.Msg) bool {
	var req models.CreateListingRequest
	if err := json.Unmarshal(msg.Data, &req); err != nil {
		log.Println("Failed to unmarshal listing data:", err)
		return false
	}

	tx, err := db.Begin()
	if err != nil {
		log.Println("DB transaction error:", err)
		return false
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRow("SELECT id FROM users WHERE email = ?", req.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		res, err := tx.Exec(`INSERT INTO users (full_name, email, phone, password, status) VALUES (?, ?, ?, '', 1)`,
			req.FullName, req.Email, req.Phone)
		if err != nil {
			log.Println("Failed to insert user:", err)
			return false
		}
		id, _ := res.LastInsertId()
		userID = int(id)
	} else if err != nil {
		log.Println("User lookup failed:", err)
		return false
	}

	point := fmt.Sprintf("POINT(%f %f)", req.Lat, req.Long)
	_, err = tx.Exec(`
		INSERT INTO listing (user_id, title, slug, category, latlong, rating, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ST_GeomFromText(?), 0, 0, ?, ?)`,
		userID, req.Title, req.Slug, req.Category, point, time.Now(), time.Now())
	if err != nil {
		log.Println("Failed to insert listing:", err)
		return false
	}

	if err := tx.Commit(); err != nil {
		log.Println("Transaction commit error:", err)
		return false
	}

	log.Println("Listing created successfully for:", req.Email)
	return true
}

/*
Usage:

1. Init() – initialize NATS & JetStream connection
2. CreateStream("MY_STREAM", "my.subject") – create stream and subject
3. ListStreams() – return list of stream names
4. ListSubjects("MY_STREAM") – return list of subjects for a stream
5. Publish("my.subject", "message content") – publish message to a subject
6. Subscribe("my.subject", handler) – realtime subscription with manual ack
7. PollMessages("my.subject", "durable-consumer", handler, 10) – poll messages every 1 minute with process function:
	func(msg *nats.Msg) bool { return true to ack, false to skip }
8. ProcessListingInsert(db, msg) – parse and insert the message data into MySQL
*/
