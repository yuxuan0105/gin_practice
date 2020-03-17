package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	User_id    string `json:"user_id"`
	Email      string `json:"email" form:"email" db:"email" binding:"required,email"`
	Password   string `json:"password" form:"password" db:"password" binding:"required"`
	Nickname   string `json:"nickname" form:"nickname" db:"nickname" binding:"required"`
	Created_on string `json:"created_on"`
}

func (this *Model) login() {

}

func (this *Model) addUser(c *gin.Context) {
	//get form
	var form User

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	//check email or nickname exist or not
	var is_exist bool
	if err := this.db.Get(
		&is_exist,
		"SELECT EXISTS(SELECT 1 FROM account WHERE email=$1 OR nickname=$2)",
		form.Email, form.Nickname,
	); err != nil {
		log.Printf("addUser db error: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	if is_exist {
		log.Printf("addUser: email or nickname exist")
		c.JSON(http.StatusBadRequest, nil)
		return
	}
	//encrypt password
	newpass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("addUser: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	form.Password = string(newpass)
	//insert data
	var ret string
	if err := namedGet(
		this.db,
		&ret,
		`INSERT INTO account(email,password,nickname) VALUES(:email,:password,:nickname) RETURNING user_id;`,
		&form,
	); err != nil {
		log.Printf("addUser: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Location", "/api/v1/users/"+ret)
	c.JSON(http.StatusCreated, nil)
}

func (this *Model) getUsers(c *gin.Context) {
	var data []User
	err := this.db.Select(&data, "SELECT user_id,email,nickname,created_on FROM account")
	if err != nil {
		log.Printf("getUsers: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, gin.H{
		"msg":  http.StatusText(http.StatusOK),
		"data": data,
	})
}

func (this *Model) getUserById(c *gin.Context) {
	uid := c.Param("uid")
	var data []User
	err := this.db.Select(&data,
		`SELECT user_id,email,nickname,created_on 
		FROM account 
		WHERE user_id = $1`, uid,
	)
	if err != nil {
		log.Printf("getUserById: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  http.StatusText(http.StatusOK),
		"data": data,
	})
}

func (this *Model) modifyUserName(c *gin.Context) {
	uid := c.Param("uid")
	newName := c.PostForm("nickname")
	if newName == "" {
		log.Printf("modifyUserName: No parameter nickname")
		c.JSON(http.StatusBadRequest, nil)
		return
	}

	if _, err := this.db.Exec(
		`UPDATE account SET nickname=$1 WHERE user_id=$2`,
		newName, uid,
	); err != nil {
		log.Printf("modifyUserName: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (this *Model) deleteUser(c *gin.Context) {
	uid := c.Param("uid")
	if _, err := this.db.Exec(
		`DELETE FROM account WHERE user_id=$1`,
		uid,
	); err != nil {
		log.Printf("deleteUser: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func namedGet(db *sqlx.DB, dest interface{}, query string, para interface{}) error {
	stmt, err := db.PrepareNamed(query)
	if err != nil {
		return err
	}
	if err := stmt.Get(dest, para); err != nil {
		return err
	}
	if err := stmt.Close(); err != nil {
		return err
	}
	return nil
}
