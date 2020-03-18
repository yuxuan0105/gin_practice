package v1

import (
	"net/http"

	"github.com/yuxuan0105/gin_practice/pkg/e"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	User_id    string `json:"user_id"`
	Email      string `json:"email"    form:"email"    db:"email"    binding:"required,email"`
	Password   string `json:"password" form:"password" db:"password" binding:"required"`
	Nickname   string `json:"nickname" form:"nickname" db:"nickname" binding:"required"`
	Created_on string `json:"created_on"`
}

func (this *Model) login() {

}

func (this *Model) register(c *gin.Context) {
	errPrefix := "addUser: "
	//get form
	var form User
	if err := c.ShouldBind(&form); err != nil {
		e.NewServErr(e.ERROR_BINDING_FORM, err).Handle(c, errPrefix)
		return
	}
	//add newuser
	var newId string
	if err := this.addNewUser(&newId, &form); err != nil {
		err.Handle(c, errPrefix)
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.Header("Location", "/api/v1/users/"+newId)
	c.JSON(http.StatusCreated, nil)
}

func (this *Model) getUsers(c *gin.Context) {
	errPrefix := "gerUsers: "
	var data []User
	err := this.db.Select(&data, "SELECT user_id,email,nickname,created_on FROM account")
	if err != nil {
		e.NewServErr(e.ERROR_WRONG_QUERY, err).Handle(c, errPrefix)
		return
	}
	c.Header("Content-Type", "application/json; charset=utf-8")
	c.JSON(http.StatusOK, gin.H{
		"msg":  http.StatusText(http.StatusOK),
		"data": data,
	})
}

func (this *Model) getUserById(c *gin.Context) {
	errPrefix := "gerUserById: "
	uid := c.Param("uid")
	var data []User
	err := this.db.Select(
		&data,
		`SELECT user_id,email,nickname,created_on 
		FROM account 
		WHERE user_id = $1`,
		uid,
	)
	if err != nil {
		e.NewServErr(e.ERROR_WRONG_QUERY, err).Handle(c, errPrefix)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  http.StatusText(http.StatusOK),
		"data": data,
	})
}

func (this *Model) modifyUserName(c *gin.Context) {
	errPrefix := "modifyUserName: "
	uid := c.Param("uid")
	newName := c.PostForm("nickname")
	if newName == "" {
		e.NewServErr(e.ERROR_BINDING_FORM, nil).Handle(c, errPrefix)
		return
	}

	if _, err := this.db.Exec(
		`UPDATE account SET nickname=$1 WHERE user_id=$2`,
		newName, uid,
	); err != nil {
		e.NewServErr(e.ERROR_WRONG_QUERY, err).Handle(c, errPrefix)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (this *Model) deleteUser(c *gin.Context) {
	errPrefix := "deleteUser: "
	uid := c.Param("uid")
	if _, err := this.db.Exec(
		`DELETE FROM account WHERE user_id=$1`,
		uid,
	); err != nil {
		e.NewServErr(e.ERROR_WRONG_QUERY, err).Handle(c, errPrefix)
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

/*if exist return error*/
func (this *Model) checkUserNotExist(form *User) *e.ServErr {
	var is_exist bool
	if err := this.db.Get(
		&is_exist,
		"SELECT EXISTS(SELECT 1 FROM account WHERE email=$1 OR nickname=$2)",
		form.Email, form.Nickname,
	); err != nil {
		return e.NewServErr(e.ERROR_WRONG_QUERY, err)
	}
	if is_exist {
		return e.NewServErr(e.ERROR_USER_EXIST, nil)
	}
	return nil
}

func (this *Model) addNewUser(userId *string, form *User) *e.ServErr {
	//check email or nickname exist or not
	var is_exist bool
	if err := this.db.Get(
		&is_exist,
		"SELECT EXISTS(SELECT 1 FROM account WHERE email=$1 OR nickname=$2)",
		form.Email, form.Nickname,
	); err != nil {
		return e.NewServErr(e.ERROR_WRONG_QUERY, err)
	}
	if is_exist {
		return e.NewServErr(e.ERROR_USER_EXIST, nil)
	}
	//encrypt password
	newpass, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		return e.NewServErr(e.ERROR_FAIL_ENCRYPT, err)
	}
	form.Password = string(newpass)
	//insert data
	if err := namedGet(
		this.db,
		userId,
		`INSERT INTO account(email,password,nickname) VALUES(:email,:password,:nickname) RETURNING user_id;`,
		&form,
	); err != nil {
		return e.NewServErr(e.ERROR_WRONG_QUERY, err)
	}

	return nil
}
