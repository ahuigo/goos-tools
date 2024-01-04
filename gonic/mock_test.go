package gonic

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

func mockContext(req *http.Request) (resp *httptest.ResponseRecorder, ctx *gin.Context) {
	resp = httptest.NewRecorder()
	ctx, _ = gin.CreateTestContext(resp)
	ctx.Request = req
	return
}