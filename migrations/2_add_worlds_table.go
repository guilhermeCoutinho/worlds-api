package main

import (
	"fmt"

	"github.com/go-pg/migrations"
)

func init() {
	err := migrations.Register(func(db migrations.DB) error {
		fmt.Println("creating table worlds")
		_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS worlds (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	user_id UUID REFERENCES users(id),
	name VARCHAR(255),
	metadata TEXT,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT now()
);
`)

		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table worlds")
		_, err := db.Exec(`DROP TABLE worlds`)
		return err
	})
	if err != nil {
		panic(err)
	}
}
