package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	CreateOn string `json:"create_on"`
}

//ADD: login

func (this *Model) getUsers(c *gin.Context) {
	//data := interface{}{}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "success",
		"data": nil,
	})
}

func (this *Model) getUserById(c *gin.Context) {

}

func (this *Model) addUser(c *gin.Context) {

}

func (this *Model) modifyUser(c *gin.Context) {

}

func (this *Model) deleteUser(c *gin.Context) {

}
