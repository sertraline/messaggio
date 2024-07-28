package services

import (
	db "github.com/sertraline/messaggio/database"
)


func GetStats() (map[string]any, error) {
	query := `
	WITH cpr AS (
		SELECT COUNT(id) AS cpt FROM messages WHERE processed = true 
	) SELECT * FROM cpr UNION SELECT COUNT(id) AS cnpt FROM messages
	`

	type stats struct{
		Cpt int `db:"cpt"`
	}
	res := make([]*stats, 0)

	rows, err := db.DBCon.Queryx(query)
	if err != nil {
		return map[string]any{
			"error": "failed to fetch stats",
		}, err
	}
	for rows.Next() {
		count := &stats{
			Cpt: 0,
		}
		err = rows.StructScan(count)
		if err != nil {
			panic(err)
		}
		res = append(res, count)
	}

	if len(res) < 1 {
		// сообщений нет
		return map[string]any{
			"Обработано сообщений": 0,
			"Всего сообщений": 0,
		}, nil
	} else if len(res) == 1 {
		// все сообщения обработаны
		return map[string]any{
			"Обработано сообщений": res[0].Cpt,
			"Всего сообщений": res[0].Cpt,
		}, nil
	} else {
		return map[string]any{
			"Обработано сообщений": res[0].Cpt,
			"Всего сообщений": res[1].Cpt,
		}, nil
	}
}
