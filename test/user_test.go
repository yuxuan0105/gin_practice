package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "github.com/yuxuan0105/gin_practice/api/v1"
)

func TestModel_getUsers(t *testing.T) {
	app, err := v1.NewModel()
	if err != nil {
		t.Fatalf("app init error: %s", err)
	}
	db := app.GetDBforTest()
	//test data
	exp := [][]interface{}{
		{"test@gmail.com", "123", "test"},
		{"test2@gmail.com", "456", "test2"},
		{"test3@gmail.com", "789", "test3"},
	}
	//defer cleanup
	defer func() {
		//cleanup
		db.MustExec("TRUNCATE account;ALTER SEQUENCE account_user_id_seq RESTART;")
	}()
	//add test data
	tx := db.MustBegin()
	defer tx.Rollback()
	for _, v := range exp {
		tx.MustExec("INSERT INTO account(email,password,nickname) VALUES ($1,$2,$3);", v...)
	}
	tx.Commit()
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
			assert.Equal(t, exp[i][0].(string), v.Email)
			assert.Equal(t, exp[i][2].(string), v.Nickname)
		}
	}
}
