package e

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

type ServErr struct {
	code int
	err  error
}

func NewServErr(code int, err error) *ServErr {
	if err == nil {
		err = fmt.Errorf("%s", GetCodeMsg(code))
	}
	return &ServErr{code, err}
}

func (this *ServErr) Error() string {
	return fmt.Sprintf("%d - %s", this.code, this.err)
}

func (this *ServErr) Wrap(prefix string) {
	this.err = fmt.Errorf("%s%w", prefix, this.err)
}

func (this *ServErr) GetCode() int {
	return this.code
}

func (this *ServErr) Handle(c *gin.Context, prefix string) {
	this.Wrap(prefix)
	statusCode := this.code / 100
	if statusCode > 499 {
		log.Println(this)
	}
	var errBody = map[string]interface{}{
		"code":    this.code,
		"message": GetCodeMsg(this.code),
	}
	c.JSON(statusCode, gin.H{
		"error": errBody,
	})
}
