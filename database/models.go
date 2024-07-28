package database

import (
	"database/sql"
	"encoding/json"
	"strings"
	"time"
)

// Обертка вокруг string для которой мы ниже определяем unmarshal
// Зачем это нужно? Чтобы различать когда пользователь намеренно прислал пустую строку и когда строка не пришла вовсе.
type NullString struct {
	sql.NullString
}

func (s *NullString) UnmarshalJSON(data []byte) error {
	s.String = strings.Trim(string(data), `"`)
	s.Valid = true
	return nil
}

func (s *NullString) MarshalJSON() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return json.Marshal(nil)
}

type Messages struct {
	Id        	int       `json:"id"`
	Content		string    `json:"content"`
	CreatedAt 	time.Time `json:"created_at" db:"created_at"`
	Processed 	bool `json:"processed"`
}
