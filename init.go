package H

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func InitGinServer() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)
	Server := gin.New()
	Server.Use(gin.Recovery())

	Server.Static("/public", "./public")

	store := cookie.NewStore([]byte("sdvsv#@R@#R$@fvsvsdvdfv"))
	Server.Use(sessions.Sessions("mysession", store))

	return Server

}
