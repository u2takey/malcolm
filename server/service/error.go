package service

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	req "github.com/arlert/malcolm/utils/reqlog"
)

var (
	ErrInvalidParam = HttpError(http.StatusBadRequest, "InvalidParameter", "Invalid parameter.")
	ErrNotFound     = HttpError(http.StatusNotFound, "NotFound", "Request resource not found.")
	ErrDB           = HttpError(http.StatusInternalServerError, "DBError", "DB error.")
	ErrInnerServErr = HttpError(http.StatusInternalServerError,
		"InternalServerError", "Internal server error, please try again later.")
)

type httpError struct {
	HttpCode   int
	ErrCode    string
	ErrMessage string
}

func HttpError(code int, errCode, errMessage string) *httpError {
	return &httpError{HttpCode: code, ErrCode: errCode, ErrMessage: errMessage}
}

func (p *httpError) WithMessage(msg string) *httpError {
	p.ErrMessage = msg
	return p
}

func (p *httpError) AppendMessage(msg string) *httpError {
	p.ErrMessage = p.ErrMessage + "\n" + msg
	return p
}

func (p *httpError) WithCode(ecode string) *httpError {
	p.ErrCode = ecode
	return p
}

func (p *httpError) Error() string {
	return fmt.Sprintf("httpcode: %d; errcode:%s; errmsg: %s", p.HttpCode, p.ErrCode, p.ErrMessage)
}

// --------------------------------------------------------------------
// Error response
func E(c *gin.Context, e *httpError) {
	req.Entry(c).WithField("error", e.Error()).WithField("path", c.Request.URL.Path).Info("")
	c.JSON(e.HttpCode, gin.H{
		"errorCode": e.ErrCode,
		"errorMsg":  e.ErrMessage,
	})
}

// Normal response
func R(c *gin.Context, body interface{}) {
	c.JSON(http.StatusOK, body)
}
