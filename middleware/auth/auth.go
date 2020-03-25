package auth

import (
	"fmt"
	"time"

	jwt "github.com/appleboy/gin-jwt"
	"github.com/gin-gonic/gin"
	v1 "github.com/yuxuan0105/gin_practice/api/v1"
	"github.com/yuxuan0105/gin_practice/middleware/database"
	"github.com/yuxuan0105/gin_practice/pkg/e"
	"golang.org/x/crypto/bcrypt"
)

func NewAuth() *jwt.GinJWTMiddleware {
	return &jwt.GinJWTMiddleware{
		Realm:         "gin_parctice",
		Key:           []byte("secret11223344"),
		Timeout:       time.Hour * 12,
		MaxRefresh:    time.Hour * 24,
		Authenticator: authenticatorFunc,
		Unauthorized:  unauthorizedFunc,
	}
}

func authenticatorFunc(c *gin.Context) (interface{}, error) {
	db := database.GetDbFromContext(c)
	var form struct {
		Email    string `form:"email"    binding:"required,email"`
		Password string `form:"password" binding:"required"`
	}
	if err := c.ShouldBind(&form); err != nil {
		return nil, jwt.ErrMissingLoginValues
	}
	var target v1.User
	if err := db.Get(&target, "SELECT * FROM account WHERE email=$1;", form.Email); err != nil {
		return nil, err
	}
	if target.User_id == "" {
		return nil, jwt.ErrFailedAuthentication
	}
	if err := bcrypt.CompareHashAndPassword([]byte(target.Password), []byte(form.Password)); err != nil {
		return nil, jwt.ErrFailedAuthentication
	}

	return &target, nil
}

func unauthorizedFunc(c *gin.Context, code int, message string) {
	e.NewErrHandler(c, "jwt: ").Handle(e.ERROR_UNAUTHORIZED, fmt.Errorf("%s", message))
}
