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
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	c "github.com/yuxuan0105/gin_practice/controller"
	"github.com/yuxuan0105/gin_practice/middleware/database"
	"github.com/yuxuan0105/gin_practice/pkg/setting"
	"golang.org/x/crypto/bcrypt"
)

var (
	app      *c.Controller
	db       *sqlx.DB
	testData [][]interface{}
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	var err error
	//init app
	app = c.NewController()
	//init db
	var v *viper.Viper
	v, err = setting.GetSetting("")
	if err != nil {
		log.Panicf("TestMain: %s", err)
	}
	//setup database
	db, err = database.SetupDatabase(v)
	if err != nil {
		log.Panicf("TestMain: %s", err)
	}
	//check table is empty
	var is_not_empty bool
	err = db.Get(&is_not_empty, "SELECT EXISTS(SELECT 1 FROM account LIMIT 1);")
	if err != nil {
		log.Panic(err)
	} else if is_not_empty {
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
		if _, err := db.Exec("TRUNCATE account;ALTER SEQUENCE account_user_id_seq RESTART;"); err != nil {
			log.Panicf("TestMain: %s", err)
		}
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

func Test_getUsers(t *testing.T) {
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

func Test_getUserById(t *testing.T) {
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

func Test_modifyUserName(t *testing.T) {
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

func Test_deleteUser(t *testing.T) {
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

func Test_addUser(t *testing.T) {
	path := "/api/v1/users"
	body := &url.Values{}
	body.Add("email", "Gary1322@gmail.com")
	body.Add("password", "123456")
	body.Add("nickname", "Gary")
	//send request
	req := newRequestWithBody("POST", path, body)
	w := app.ServeTestRequest(req)
	//check code
	if !assert.Equal(t, http.StatusCreated, w.Code) {
		t.Fatal()
	}
	assert.Equal(t, "/api/v1/users/4", w.Header().Get("Location"))
	//check add success or not
	var ret []string
	if err := db.Select(
		&ret,
		"SELECT password FROM account WHERE email=$1 AND nickname=$2;",
		body.Get("email"), body.Get("nickname"),
	); err != nil {
		t.Fatal()
	}
	if assert.Equal(t, 1, len(ret)) {
		if err := bcrypt.CompareHashAndPassword([]byte(ret[0]), []byte(body.Get("password"))); err != nil {
			t.Fail()
		}
	}
	//send the same request again
	w = app.ServeTestRequest(req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	//test no parameter request
	body = &url.Values{}
	req = newRequestWithBody("POST", path, body)
	w = app.ServeTestRequest(req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}
