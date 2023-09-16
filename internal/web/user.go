package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/huangyul/go-server/internal/domain"
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
	server.POST("/user/login", u.LoginJWT)
	server.POST("/user/profile", u.Profile)
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

	err = u.srv.SignUp(ctx, &domain.User{Password: req.Password, Email: req.Email})
	if errors.Is(err, errors.New("邮箱冲突")) {
		ctx.String(http.StatusInternalServerError, "邮箱冲突")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}

	ctx.String(http.StatusOK, "注册成功")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	err := ctx.Bind(&req)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	uDomian := domain.User{
		Password: req.Password,
		Email:    req.Email,
	}
	uDomain, err := u.srv.FindByEmail(ctx, uDomian)
	if errors.Is(err, service.ErrNotFound) {
		ctx.String(http.StatusOK, "user not exist")
		return
	}
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "Incorrect email or passwrod")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	sess := sessions.Default(ctx)
	sess.Set("userId", uDomain.ID)
	sess.Options(sessions.Options{
		MaxAge: 60,
	})
	sess.Save()

	ctx.String(http.StatusOK, "login successful")
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type ReqParam struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	var reqParam ReqParam
	err := ctx.Bind(&reqParam)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "system error")
		return
	}
	if reqParam.Email == "" {
		ctx.String(http.StatusBadRequest, "email cannot be empty")
		return
	}
	if reqParam.Password == "" {
		ctx.String(http.StatusBadRequest, "password cannot be empty")
		return
	}
	du := domain.User{
		Password: reqParam.Password,
		Email:    reqParam.Email,
	}
	du, err = h.srv.FindByEmail(ctx, du)
	if err == service.ErrNotFound {
		ctx.String(http.StatusOK, "user not found")
		return
	}
	if err == service.ErrInvalidUserOrPassword {
		ctx.String(http.StatusOK, "Incorrect emial addr or password")
		return
	}
	if err != nil {
		ctx.String(http.StatusInternalServerError, "system error")
	}
	claim := UserClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
		UId:       du.ID,
		UserAgent: ctx.Request.UserAgent(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claim)
	tokenStr, err := token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "system error")
		return
	}
	ctx.String(http.StatusOK, tokenStr)
}

func (h *UserHandler) Profile(ctx *gin.Context) {}

type UserClaim struct {
	jwt.RegisteredClaims
	UId       int64
	UserAgent string
}
