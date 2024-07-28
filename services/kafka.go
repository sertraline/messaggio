package services

import (
	"context"

	"github.com/segmentio/kafka-go"
	db "github.com/sertraline/messaggio/database"
)


func GetKafkaMessageFromTopic(ctx context.Context) (map[string]any, error) {
	r := ctx.Value(db.CtxKey).(*kafka.Reader)

    m, err := r.ReadMessage(context.Background())
    if err != nil {
        return map[string]any{
			"key": "",
			"value": "",
			"topic": "",
			"partition": "",
			"offset": "",
		}, err
    }

	return map[string]any{
		"key": string(m.Key),
		"value": string(m.Value),
		"topic": string(m.Topic),
		"partition": m.Partition,
		"offset": m.Offset,
		}, nil
}
