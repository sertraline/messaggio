package validators

import (
	"errors"
	"net/http"

	models "github.com/sertraline/messaggio/database"
)

type MessageSaveRequest struct {
	Content models.NullString `json:"content"`
}

func (u *MessageSaveRequest) Bind(r *http.Request) error {
	if !u.Content.Valid || u.Content.String == "" {
		return errors.New("content is missing")
	}

	return nil
}
