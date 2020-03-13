package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	v1 "github.com/yuxuan0105/gin_practice/api/v1"
)

var (
	app      *v1.Model
	db       *sqlx.DB
	testData [][]interface{}
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	var err error
	//init app
	app, err = v1.NewModel()
	if err != nil {
		log.Panicf("app init error: %s", err)
	}
	db = app.GetDBforTest()
	//check table is empty
	var is_not_empty int
	err = db.Get(&is_not_empty, "SELECT CASE WHEN EXISTS (SELECT * FROM account LIMIT 1) THEN 1 ELSE 0 END;")
	if err != nil {
		log.Panic(err)
	} else if is_not_empty == 1 {
		log.Panic("table is not empty")
	}
	//test data
	testData = [][]interface{}{
		{"test@gmail.com", "123", "test"},
		{"test2@gmail.com", "456", "test2"},
		{"test3@gmail.com", "789", "test3"},
	}
	//defer cleanup
	defer func() {
		db.MustExec("TRUNCATE account;ALTER SEQUENCE account_user_id_seq RESTART;")
	}()
	//add test data
	tx := db.MustBegin()
	defer tx.Rollback()
	for _, v := range testData {
		tx.MustExec("INSERT INTO account(email,password,nickname) VALUES ($1,$2,$3);", v...)
	}
	tx.Commit()
	//Run tests
	m.Run()
}

func TestModel_getUsers(t *testing.T) {
	//send request
	w := app.ServeTestRequest("GET", "/api/v1/users")
	//check code
	if !assert.Equal(t, http.StatusOK, w.Code) {
		t.Fatal()
	}
	//should have three object in res["data"]
	var res testBody
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("error at unmarshal json: %s", err)
	}
	if assert.Equal(t, 3, len(res.Data)) {
		for i, v := range res.Data {
			assert.Equal(t, testData[i][0].(string), v.Email)
			assert.Equal(t, testData[i][2].(string), v.Nickname)
		}
	}
}

func TestModel_getUserById(t *testing.T) {
	testId := 2
	url := fmt.Sprintf("/api/v1/users/%d", testId)
	//send request
	w := app.ServeTestRequest("GET", url)
	//check code
	if !assert.Equal(t, http.StatusOK, w.Code) {
		t.Fatal()
	}
	//should have three object in res["data"]
	var res testBody
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("error at unmarshal json: %s", err)
	}
	if assert.Equal(t, 1, len(res.Data)) {
		assert.Equal(t, testData[testId-1][0].(string), res.Data[0].Email)
		assert.Equal(t, testData[testId-1][2].(string), res.Data[0].Nickname)
	}
}
