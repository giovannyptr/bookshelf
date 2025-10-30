package api

import "github.com/gin-gonic/gin"

type Resp struct {
	OK      bool        `json:"ok"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func OK(c *gin.Context, data any)               { c.JSON(200, Resp{OK: true, Data: data}) }
func Created(c *gin.Context, data any)          { c.JSON(201, Resp{OK: true, Data: data}) }
func Fail(c *gin.Context, code int, msg string) { c.JSON(code, Resp{OK: false, Error: msg}) }
