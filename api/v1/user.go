package v1

import "github.com/gin-gonic/gin"

type User struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	CreateOn string `json:"create_on"`
}

//ADD: login

func (this *Model) getUsers(c *gin.Context) {

}

func (this *Model) getUserById(c *gin.Context) {

}

func (this *Model) addUser(c *gin.Context) {

}

func (this *Model) modifyUser(c *gin.Context) {

}

func (this *Model) deleteUser(c *gin.Context) {

}
