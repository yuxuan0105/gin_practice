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

func newUser() *User {
	return &User{}
}

//ADD: login

func (this *Model) getUsers(c *gin.Context) {
	var data []User
	err := this.db.Select(&data, "SELECT user_id,email,nickname,created_on FROM account")
	if err != nil {
		log.Printf("getUsers: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg":  http.StatusText(http.StatusInternalServerError),
			"data": nil,
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  http.StatusText(http.StatusOK),
		"data": data,
	})
}

func (this *Model) getUserById(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"msg":  http.StatusText(http.StatusOK),
		"data": nil,
	})
}

func (this *Model) addUser(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{
		"msg":  http.StatusText(http.StatusCreated),
		"data": nil,
	})
}

func (this *Model) modifyUser(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{
		"msg":  http.StatusText(http.StatusNoContent),
		"data": nil,
	})
}

func (this *Model) deleteUser(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{
		"msg":  http.StatusText(http.StatusNoContent),
		"data": nil,
	})
}
