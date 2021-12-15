package oauth

import (
	"net/http"
	"os"
	"sync"

	"github.com/choby/oauth2.0_golang/server/libs/request"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-session/session"
)

var (
	srv  *server.Server
	once sync.Once
)

// 初始化服务
func InitServer(manager *manage.Manager) *server.Server {
	once.Do(func() {
		// oauth2Server = server.NewServer(server.NewConfig(), manager)
		srv = server.NewDefaultServer(manager)
		SetAllowedGrantType("authorization_code", "refresh_token")
		SetAllowGetAccessRequest(true)

		// 密码授权模式才需要用到这个配置, 这个模式不需要分配授权码,而是直接分配token,通常用于无后端的应用
		SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
			if username == "test" && password == "test" {
				userID = "test"
			}
			return userID, nil
		})

		// 这一行很关键,这个方法让oauth框架识别当前用户身份标识(并且可以人为处理登陆状态检验等等)
		// 具体看userAuthorizeHandler方法实现
		SetUserAuthorizationHandler(userAuthorizeHandler)
	})

	return srv
}

// 处理身份认证请求, the authorization request handling
func HandleAuthorizeRequest(c *gin.Context) {
	err := srv.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Abort()
}

// 处理token请求
func HandleTokenRequest(c *gin.Context) {
	err := srv.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(500, err)
		return
	}
	c.Abort()
}

// oauth框架通过本方法识别用户身份信息,并且可以人为进行登录状态校验
// 本方法正常执行后,则会为客户端分配授权码(authorization_code)
func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	_ = request.DumpRequest(os.Stdout, "userAuthorizeHandler", r) // Ignore the error

	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		store.Set("ReturnUri", r.Form)
		store.Save()

		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()
	return
}
