package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	db "github.com/sertraline/messaggio/database"
	errResponse "github.com/sertraline/messaggio/errors"
	validators "github.com/sertraline/messaggio/validators"
)


func SaveMessage(data *validators.MessageSaveRequest, ctx context.Context) (int, error) {
	w := &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092", "localhost:9093", "localhost:9094"),
		Topic:   "messages",
		Balancer: &kafka.LeastBytes{},
	}
	query := ""
	uniq := uuid.New()
	message := kafka.Message{
        Key:   []byte(uniq.String()),
        Value: []byte(data.Content.String),
    }

	err := w.WriteMessages(context.Background(),
		message,
	)
	if err != nil {
		// в случае ошибки записи сообщение попадает в БД необработанным
		fmt.Println("failed to write messages:", err)
		query = `
		INSERT INTO Messages (content) VALUES ($1);
		`
	} else {
		// в случае успеха processed = true
		query = `
		INSERT INTO Messages (content, processed) VALUES ($1, true);
		`
	}

	if err := w.Close(); err != nil {
		fmt.Println("failed to close writer:", err)
	}

	_, err = db.DBCon.Exec(query, data.Content.String)
	if err != nil {
		fmt.Println(err)
		return 1, err
	}

	return 0, nil
}

func GetMessageByID(id int) (*db.Messages, error) {
	query := `
		SELECT * FROM Messages WHERE id = $1;
	`

	m := &db.Messages{ Id: 0 }
	err := db.DBCon.Get(m, query, id)
	if err == sql.ErrNoRows {
		fmt.Println(err)
		return m, errResponse.ErrNotExists
	} else if err != nil {
		fmt.Println(err)
		return m, errResponse.ErrDB
	}

	return m, nil
}
