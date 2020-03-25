package database

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	cst "github.com/yuxuan0105/gin_practice/pkg/constant"
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

func GetDbFromContext(c *gin.Context) *sqlx.DB {
	db, ok := c.MustGet(cst.DB_KEY).(*sqlx.DB)
	if !ok {
		log.Panicln("GetDbFromContext: No such value")
	}
	return db
}

func GetMiddlewareFunc(db *sqlx.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(cst.DB_KEY, db)
		c.Next()
	}
}
