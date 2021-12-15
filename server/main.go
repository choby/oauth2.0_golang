package main

import (
	"net/http"

	"github.com/choby/oauth2.0_golang/server/oauth"
	"github.com/choby/oauth2.0_golang/server/routes"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/store"
)

func main() {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token存储方式
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	// manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	manager.MapAccessGenerate(generates.NewAccessGenerate())

	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
	})
	manager.MapClientStorage(clientStore)

	// 初始化oauth2.0服务
	oauth.InitServer(manager)

	g := gin.Default()

	routes.SetAuthRoutes(g)

	api := g.Group("/api")
	{
		api.Use(oauth.HandleTokenRequest)
		// api.Use(func(c *gin.Context) {
		// 	token, err := oauth.Oauth2Server.ValidationBearerToken(c.Request)
		// 	if err != nil {
		// 		c.AbortWithError(http.StatusBadRequest, err)
		// 		// http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		// 		return
		// 	}
		// 	_, err1 := mgr.GetClient(c, token.GetClientID())
		// 	if err1 != nil {
		// 		c.AbortWithError(http.StatusBadRequest, err1)
		// 		return
		// 	}
		// })
		api.GET("/test", func(c *gin.Context) {
			ti, exists := c.Get(oauth.DefaultConfig.TokenKey)
			if exists {
				c.JSON(http.StatusOK, ti)
				return
			}
			c.String(http.StatusOK, "no found")
		})
	}

	g.Run(":9096")
}
