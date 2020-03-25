package e

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type ErrHandler struct {
	c      *gin.Context
	prefix string
}

func NewErrHandler(c *gin.Context, prefix string) *ErrHandler {
	return &ErrHandler{c, prefix}
}

func (this *ErrHandler) Handle(code int, err error) {
	statusCode := code / 100
	if err == nil {
		err = errors.New(GetCodeMsg(code))
	}
	if statusCode > 499 {
		log.Printf("%d - %s%s", code, this.prefix, err)
	}
	var errBody = map[string]interface{}{
		"code":    code,
		"message": fmt.Sprintln(err),
	}
	this.c.JSON(statusCode, gin.H{
		"error": errBody,
	})
}
