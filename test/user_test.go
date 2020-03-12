package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	v1 "github.com/yuxuan0105/gin_practice/api/v1"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testBody struct {
	Msg  string    `json:"msg"`
	Data []v1.User `json:"data"`
}

func TestModel_getUsers(t *testing.T) {
	app, err := v1.NewModel()
	if err != nil {
		t.Fatalf("app init error: %s", err)
	}
	//get tx
	tx, err := app.GetTxForTest()
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		//clean up
		_, e := tx.Exec("TRUNCATE account;ALTER SEQUENCE account_user_id_seq RESTART;")
		if e != nil {
			t.Logf("error at cleanup table: %s", e)
		}

		if e := tx.Commit(); e != nil {
			tx.Rollback()
		}
	}()
	//add three rows
	exp := [][]string{
		{"test@gmail.com", "123", "test"},
		{"test2@gmail.com", "456", "test2"},
		{"test3@gmail.com", "789", "test3"},
	}

	for _, v := range exp {
		_, err := tx.Exec("INSERT INTO account (email,password,nickname) VALUES ($1,$2,$3);", v[0], v[1], v[2])
		if err != nil {
			t.Logf("error at insert query: %s", err)
		}
	}
	//send request
	w := app.ServeTestRequest("GET", "/api/v1/users")
	var res testBody
	err = json.Unmarshal(w.Body.Bytes(), &res)
	if err != nil {
		t.Fatalf("error at unmarshal json: %s", err)
	}
	//should have three object in res["data"]
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", res.Msg)
	if assert.Equal(t, 3, len(res.Data)) {
		for i, v := range res.Data {
			assert.Equal(t, exp[i][0], v.Email)
			assert.Equal(t, exp[i][2], v.Nickname)
		}
	}
}
