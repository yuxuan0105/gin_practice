package test

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/yuxuan0105/gin_practice/api/v1"
)

func init() {
	gin.SetMode(gin.TestMode)
}

type testBody struct {
	Msg  string    `json:"msg"`
	Data []v1.User `json:"data"`
}
