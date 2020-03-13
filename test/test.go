package test

import (
	v1 "github.com/yuxuan0105/gin_practice/api/v1"
)

type testBody struct {
	Msg  string    `json:"msg"`
	Data []v1.User `json:"data"`
}
