package v1

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type User struct {
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	CreateOn string `json:"create_on"`
}

func newUser() *User {
	return &User{}
}

//ADD: login

func (this *Model) getUsers(c *gin.Context) {
	data, err := getUserDatas(this.db, "")
	if err != nil {
		log.Println(err)
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

func getUserDatas(db *sql.DB, condition string) ([]*User, error) {
	query := "SELECT email, nickname,created_on From account"
	if condition != "" {
		query += " WHERE " + condition
	}
	query += ";"

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error at querying data: %s", err)
	}
	defer rows.Close()

	data := []*User{}
	for rows.Next() {
		temp := newUser()
		if err := rows.Scan(&temp.Email, &temp.Nickname, &temp.CreateOn); err != nil {
			return nil, fmt.Errorf("error at scaning data: %s", err)
		}
		data = append(data, temp)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return data, nil
}
