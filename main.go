package main

import (
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-server/internal/repository"
	"github.com/huangyul/go-server/internal/repository/dao"
	"github.com/huangyul/go-server/internal/service"
	"github.com/huangyul/go-server/internal/web"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func initWeb() *gin.Engine {
	server := gin.Default()

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	uDao := dao.NewUserDAO(db)
	uRep := repository.NewUserRepository(uDao)
	uSrv := service.NewUserService(uRep)
	user := web.NewUserHandler(uSrv)

	err := dao.InitUser(db)
	if err != nil {
		panic(err)
	}

	return user
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:abc123456@tcp(47.106.214.127:3306)/basic"))
	if err != nil {
		panic(err)
	}
	return db
}

func main() {
	db := initDB()
	server := initWeb()

	user := initUser(db)

	server.Use(cors.New(cors.Config{
		AllowHeaders: []string{"Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			return strings.HasPrefix(origin, "http://localhost")
		},
	}))

	user.RegisterRoutes(server)

	server.Run(":8882")
}
