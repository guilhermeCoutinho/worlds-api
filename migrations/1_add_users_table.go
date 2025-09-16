package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	err := migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating table users")
		_, err := db.Exec(`
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
`)

		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table users")
		_, err := db.Exec(`DROP TABLE users`)
		return err
	})
	if err != nil {
		panic(err)
	}
}
