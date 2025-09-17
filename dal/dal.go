package dal

import (
	"github.com/go-pg/pg"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type DAL struct {
	db                    *pg.DB
	WorldsDAL             WorldsDAL
	UserDAL               UserDAL
	WorldsTransferJobsDAL WorldsTransferJobsDAL
}

func ConnectDB(config *viper.Viper) *pg.DB {
	return pg.Connect(&pg.Options{
		User:     config.GetString("postgres.user"),
		Password: config.GetString("postgres.password"),
		Addr:     config.GetString("postgres.addr"),
		Database: config.GetString("postgres.database"),
	})
}

func NewDAL(db *pg.DB, redisClient *redis.Client) *DAL {
	return &DAL{
		db:                    db,
		WorldsDAL:             NewWorldsDAL(db, redisClient),
		UserDAL:               NewUserDAL(db),
		WorldsTransferJobsDAL: NewWorldsTransferJobsDAL(db),
	}
}
