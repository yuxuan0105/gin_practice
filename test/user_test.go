package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "github.com/yuxuan0105/gin_practice/api/v1"
	"github.com/yuxuan0105/gin_practice/database"
)

func TestModel_getUsers(t *testing.T) {
	app, err := v1.NewModel()
	if err != nil {
		t.Fatalf("app init error: %s", err)
	}

	//add three rows
	exp := [][]interface{}{
		{"test@gmail.com", "123", "test"},
		{"test2@gmail.com", "456", "test2"},
		{"test3@gmail.com", "789", "test3"},
	}

	db, _ := app.GetDBforTest()
	if err := database.InsertDatas(db, "account(email,password,nickname)", exp); err != nil {
		t.Fatal(err)
	}
	//send request
	w := app.ServeTestRequest("GET", "/api/v1/users")
	var res testBody
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("error at unmarshal json: %s", err)
	}
	//should have three object in res["data"]
	assert.Equal(t, http.StatusOK, w.Code)
	if assert.Equal(t, 3, len(res.Data)) {
		for i, v := range res.Data {
			assert.Equal(t, exp[i][0].(string), v.Email)
			assert.Equal(t, exp[i][2].(string), v.Nickname)
		}
	}
	if err := database.CleanupTable(db, "account", "user_id"); err != nil {
		t.Log(err)
	}
}
