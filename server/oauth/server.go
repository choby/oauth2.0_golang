package oauth

import (
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
)

var (
	oauth2Server *server.Server
	once         sync.Once
)

// 初始化服务
func InitServer(manager *manage.Manager) *server.Server {
	once.Do(func() {
		// oauth2Server = server.NewServer(server.NewConfig(), manager)
		oauth2Server = server.NewDefaultServer(manager)
	})

	return oauth2Server
}

// 处理身份认证请求, the authorization request handling
func HandleAuthorizeRequest(c *gin.Context) {
	err := oauth2Server.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Abort()
}

// 处理token请求
func HandleTokenRequest(c *gin.Context) {
	err := oauth2Server.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Abort()
}
