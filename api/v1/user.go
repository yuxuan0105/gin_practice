package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	User_id    string `json:"user_id"`
	Email      string `json:"email"`
	Nickname   string `json:"nickname"`
	Created_on string `json:"created_on"`
}

//ADD: login

func (this *Model) getUsers(c *gin.Context) {
	var data []User
	err := this.db.Select(&data, "SELECT user_id,email,nickname,created_on FROM account")
	if err != nil {
		log.Printf("getUsers: %s", err)
		c.JSON(http.StatusInternalServerError, nil)
		return
	}
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

func (this *Model) addUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"msg":  http.StatusText(http.StatusCreated),
		"data": nil,
	})
}

func (this *Model) modifyUserName(c *gin.Context) {
	uid := c.Param("uid")
	newName := c.PostForm("nickname")
	if newName == "" {
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
