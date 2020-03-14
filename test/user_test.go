package test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
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
	req := newRequest("GET", "/api/v1/users")
	w := app.ServeTestRequest(req)
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
			temp := v.(map[string]interface{})
			assert.Equal(t, testData[i][0].(string), temp["email"])
			assert.Equal(t, testData[i][2].(string), temp["nickname"])
		}
	}
}

func TestModel_getUserById(t *testing.T) {
	testId := 1
	path := fmt.Sprintf("/api/v1/users/%d", testId)
	//send request
	req := newRequest("GET", path)
	w := app.ServeTestRequest(req)
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
		temp := res.Data[0].(map[string]interface{})
		assert.Equal(t, testData[testId-1][0].(string), temp["email"])
		assert.Equal(t, testData[testId-1][2].(string), temp["nickname"])
	}
}

func TestModel_modifyUserName(t *testing.T) {
	testId := 1
	path := fmt.Sprintf("/api/v1/users/%d", testId)
	//send request
	body := &url.Values{}
	body.Add("nickname", "John")
	req := newRequestWithBody("PATCH", path, body)
	w := app.ServeTestRequest(req)
	//check code
	if !assert.Equal(t, http.StatusNoContent, w.Code) {
		t.Fatal()
	}
	//check modify success or not
	var target string
	db.Get(&target, "SELECT nickname FROM account WHERE user_id=$1;", testId)
	assert.Equal(t, body.Get("nickname"), target)
	//test no parameter request
	body = &url.Values{}
	req = newRequestWithBody("PATCH", path, body)
	w = app.ServeTestRequest(req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestModel_deleteUser(t *testing.T) {
	testId := 3
	path := fmt.Sprintf("/api/v1/users/%d", testId)
	//send request
	req := newRequest("DELETE", path)
	w := app.ServeTestRequest(req)
	//check code
	if !assert.Equal(t, http.StatusNoContent, w.Code) {
		t.Fatal()
	}

	var is_exist int
	if err := db.Get(
		&is_exist,
		"SELECT CASE WHEN EXISTS (SELECT * FROM account WHERE user_id=$1) THEN 1 ELSE 0 END;",
		testId,
	); err != nil {
		t.Fatal()
	}
	assert.Equal(t, 0, is_exist)
}

func TestModel_addUser(t *testing.T) {
	//send request
	body := &url.Values{}
	body.Add("email", "John132@gmail.com")
	body.Add("password", "123456")
	body.Add("nickname", "John")
	req := newRequestWithBody("POST", "/api/v1/users", body)
	w := app.ServeTestRequest(req)
	//check code
	if !assert.Equal(t, http.StatusNoContent, w.Code) {
		t.Fatal()
	}
	//check add success or not
	var is_exist int
	if err := db.Get(
		&is_exist,
		"SELECT CASE WHEN EXISTS (SELECT * FROM account WHERE email=$1 AND password=$2 AND nickname=$3) THEN 1 ELSE 0 END;",
		body.Get("email"), body.Get("password"), body.Get("nickname"),
	); err != nil {
		t.Fatal()
	}
	assert.Equal(t, 1, is_exist)
}
