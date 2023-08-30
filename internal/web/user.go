package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-server/internal/service"
)

type UserHandler struct {
	srv *service.UserService
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	return &UserHandler{
		srv: srv,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.GET("/user/signup", u.Signup)
	server.GET("/user/login", u.Login)
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	ctx.String(http.StatusOK, "登录成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {}
