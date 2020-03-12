package database

import (
	"database/sql"
	"fmt"
	"strconv"
)

func InsertDatas(db *sql.DB, tableAndColumn string, data [][]interface{}) error {
	//tx begin
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error at tx begin: %s", err)
	}

	defer func() {
		//tx commit
		if e := tx.Commit(); e != nil {
			tx.Rollback()
		}
	}()

	que := "INSERT INTO " + tableAndColumn + " VALUES ("
	for i := range data[0] {
		if i > 0 {
			que += ","
		}
		que += "$" + strconv.Itoa(i+1)
	}
	que += ");"

	for _, v := range data {
		_, err := tx.Exec(que, v...)
		if err != nil {
			return fmt.Errorf("error at insert query: %s", err)
		}
	}
	return nil
}

func CleanupTable(db *sql.DB, tableName, idName string) error {
	query := "TRUNCATE " + tableName + ";ALTER SEQUENCE " + tableName + "_" + idName + "_seq RESTART;"
	_, err := db.Exec(query)
	return err
}
