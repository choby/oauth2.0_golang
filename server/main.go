package main

import (
	"github.com/choby/oauth2.0_golang/server/oauth"
	"github.com/choby/oauth2.0_golang/server/routes"
	"github.com/gin-gonic/gin"
)

func main() {

	// 初始化oauth2.0服务
	oauth.InitServer()

	g := gin.Default()

	routes.SetAuthRoutes(g)

	g.Run(":9096")
}
