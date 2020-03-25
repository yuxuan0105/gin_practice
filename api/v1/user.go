package v1

import (
	"net/http"

	d "github.com/yuxuan0105/gin_practice/middleware/database"
	"github.com/yuxuan0105/gin_practice/pkg/e"
	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

type User struct {
	User_id    string `json:"user_id"`
	Email      string `json:"email"    db:"email"`
	Password   string `json:"password" db:"password"`
	Nickname   string `json:"nickname" db:"nickname"`
	Created_on string `json:"created_on"`
}

func SignUp(c *gin.Context) {
	errHan := e.NewErrHandler(c, "SignUp: ")
	db := d.GetDbFromContext(c)
	var form struct {
		Email    string `form:"email"    binding:"required,email"`
		Password string `form:"password" binding:"required"`
		Nickname string `form:"nickname" binding:"required"`
	}
	//get form
	if err := c.ShouldBind(&form); err != nil {
		errHan.Handle(e.ERROR_BINDING_FORM, err)
		return
	}
	//check email or nickname exist or not
	var is_exist bool
	if err := db.Get(
		&is_exist,
		"SELECT EXISTS(SELECT 1 FROM account WHERE email=$1 OR nickname=$2)",
		form.Email, form.Nickname,
	); err != nil {
		errHan.Handle(e.ERROR_WRONG_QUERY, err)
	}
	if is_exist {
		errHan.Handle(e.ERROR_USER_EXIST, nil)
	}
	//encrypt password
	newpass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		errHan.Handle(e.ERROR_FAIL_ENCRYPT, err)
	}
	form.Password = string(newpass)
	var newId string
	//insert data
	if err := db.Get(
		&newId,
		`INSERT INTO account(email,password,nickname) VALUES($1,$2,$3) RETURNING user_id;`,
		form.Email, form.Password, form.Nickname,
	); err != nil {
		errHan.Handle(e.ERROR_WRONG_QUERY, err)
	}
	//
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Location", "/api/v1/users/"+newId)
	c.JSON(http.StatusCreated, nil)
}

func GetUsers(c *gin.Context) {
	errHan := e.NewErrHandler(c, "GetUsers: ")
	db := d.GetDbFromContext(c)
	var data []User
	err := db.Select(&data, "SELECT user_id,email,nickname,created_on FROM account")
	if err != nil {
		errHan.Handle(e.ERROR_WRONG_QUERY, err)
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, gin.H{
		"msg":  http.StatusText(http.StatusOK),
		"data": data,
	})
}

func GetUserById(c *gin.Context) {
	errHan := e.NewErrHandler(c, "GetUserById: ")
	db := d.GetDbFromContext(c)
	uid := c.Param("uid")
	var data []User
	err := db.Select(
		&data,
		`SELECT user_id,email,nickname,created_on 
		FROM account 
		WHERE user_id = $1`,
		uid,
	)
	if err != nil {
		errHan.Handle(e.ERROR_WRONG_QUERY, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  http.StatusText(http.StatusOK),
		"data": data,
	})

}

func ModifyUserName(c *gin.Context) {
	errHan := e.NewErrHandler(c, "ModifyUserName: ")
	db := d.GetDbFromContext(c)
	uid := c.Param("uid")
	newName := c.PostForm("nickname")
	if newName == "" {
		errHan.Handle(e.ERROR_BINDING_FORM, nil)
		return
	}

	if _, err := db.Exec(
		`UPDATE account SET nickname=$1 WHERE user_id=$2`,
		newName, uid,
	); err != nil {
		errHan.Handle(e.ERROR_WRONG_QUERY, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func DeleteUser(c *gin.Context) {
	errHan := e.NewErrHandler(c, "DeleteUser: ")
	db := d.GetDbFromContext(c)
	uid := c.Param("uid")
	if _, err := db.Exec(
		`DELETE FROM account WHERE user_id=$1`,
		uid,
	); err != nil {
		errHan.Handle(e.ERROR_WRONG_QUERY, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
