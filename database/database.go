package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func SetupDatabase(v *viper.Viper) (*sqlx.DB, error) {
	dbConf := v.GetStringMap("database")
	connStr := fmt.Sprintf(
		"user=%s password=%s host=%s port=%d dbname=%s sslmode=disable",
		dbConf["user"], dbConf["password"], dbConf["host"], dbConf["port"], dbConf["dbname"],
	)
	newdb, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error at setuping database: %s", err)
	}
	return newdb, nil
}
