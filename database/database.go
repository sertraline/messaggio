package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DBCon *sqlx.DB

var schema = `
CREATE TABLE IF NOT EXISTS messages (
	id		     SERIAL PRIMARY KEY,
	content		 TEXT,
	created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
	processed	 BOOLEAN NOT NULL DEFAULT false
);
`

func InitDatabase() {
	var err error
	DBCon, err = sqlx.Connect("postgres", "user=postgres dbname=messaggio password=123 sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to postgresql")

	DBCon.MustExec(schema)
}
