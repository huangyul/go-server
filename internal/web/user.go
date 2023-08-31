package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-server/internal/service"

	regexp "github.com/dlclark/regexp2"
)

type UserHandler struct {
	srv         *service.UserService
	emailReg    *regexp.Regexp
	passwordReg *regexp.Regexp
}

func NewUserHandler(srv *service.UserService) *UserHandler {
	emailRegexp := `^[\w.-]+@[a-zA-Z\d]+.[a-zA-Z]{2,}$`
	passworRegexp := `^(?=.*[a-zA-Z])(?=.*[0-9])[a-zA-Z\d]{6,18}$`

	return &UserHandler{
		srv:         srv,
		emailReg:    regexp.MustCompile(emailRegexp, regexp.None),
		passwordReg: regexp.MustCompile(passworRegexp, regexp.None),
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/user/signup", u.Signup)
	server.GET("/user/login", u.Login)
}

func (u *UserHandler) Signup(ctx *gin.Context) {
	type SignupReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req SignupReq
	err := ctx.Bind(&req)
	if err != nil {
		return
	}
	ok, err := u.emailReg.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusBadRequest, "服务器出错")
		return
	}
	if !ok {
		ctx.String(http.StatusBadRequest, "please enter the correct email")
		return
	}
	ok, err = u.passwordReg.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusBadRequest, "服务器出错")
		return
	}
	if !ok {
		ctx.String(http.StatusBadRequest, "Password must use a combination of letters and numbers and be between 6 and 18 length")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusBadRequest, "the password is not the same for both entries")
		return
	}

	fmt.Printf("%v", req)
	ctx.JSON(http.StatusOK, req)
}

func (u *UserHandler) Login(ctx *gin.Context) {}
