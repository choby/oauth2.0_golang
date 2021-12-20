package routes

import (
	"net/http"

	"github.com/choby/oauth2.0_golang/server/controllers"
	"github.com/choby/oauth2.0_golang/server/oauth"
	"github.com/gin-gonic/gin"
)

func SetAuthRoutes(g *gin.Engine) {
	// 加载静态文件
	g.LoadHTMLGlob("static/*")
	// 渲染登录页
	g.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "login",
		})
	})

	g.POST("/login", controllers.Login) // 登录
	g.GET("/auth", controllers.Auth)    // 授权页面,选择需要授权的权限项
	g.POST("/authorize", controllers.Authorize)
	g.POST("/token", oauth.HandleTokenRequest) // 应用程序通过此请求获取token

}
