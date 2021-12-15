package main

import (
	"io"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/choby/oauth2.0_golang/server/oauth"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
)

var (
	dumpvar = true
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
	clientStore.Set("22222", &models.Client{
		ID:     "22222",
		Secret: "22222",
		Domain: "http://localhost:9094",
	})
	manager.MapClientStorage(clientStore)

	// 初始化oauth2.0服务
	oauth.InitServer(manager)
	// oauth.SetAllowedGrantType("authorization_code", "refresh_token")
	oauth.SetAllowGetAccessRequest(true)
	oauth.SetClientInfoHandler(server.ClientFormHandler)

	// 密码授权模式才需要用到这个配置, 这个模式不需要分配授权码,而是直接分配token,通常用于无后端的应用
	// oauth.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
	// 	if username == "test" && password == "test" {
	// 		userID = "test"
	// 	}
	// 	return userID, nil
	// })

	// 这一行很关键,这个方法让oauth框架识别当前用户身份标识(并且可以认为处理登陆状态检验等等)
	// 具体看userAuthorizeHandler方法实现
	oauth.SetUserAuthorizationHandler(userAuthorizeHandler)

	g := gin.Default()

	// 加载静态文件
	g.LoadHTMLGlob("static/*")
	// 渲染登录页
	g.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "login",
		})
	})
	// 登录
	g.POST("/login", func(c *gin.Context) {
		// fmt.Printf("%v \n", c.PostForm("username"))
		// fmt.Printf("%v", c.Request.PostForm["username"][0])
		// if c.PostForm("usename") == nil {
		// 	c.String(http.StatusInternalServerError, "username is empty")
		// 	return
		// }
		// if c.Request.PostForm["password"] == nil {
		// 	c.String(http.StatusInternalServerError, "password is empty")
		// 	return
		// }
		var form LoginForm
		if c.ShouldBind(&form) == nil {

			store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
			if err != nil {
				c.String(http.StatusInternalServerError, "session start error")
				return
			}
			if form.UserName == "test" && form.Password == "test" {
				store.Set("LoggedInUserID", "test")
				store.Save()
				// c.Writer.Header().Set("Location", "/auth")
				// c.Writer.WriteHeader(http.StatusFound)
				// 登录成功,跳转至授权页面,选择需要授权的权限项
				c.Redirect(http.StatusFound, "/auth")
			} else {
				c.String(http.StatusInternalServerError, "username or password error")
			}
		}

	})

	// 授权页面,选择需要授权的权限项
	g.GET("/auth", authHandler)
	g.POST("/oauth/authorize", func(c *gin.Context) {
		if dumpvar {
			dumpRequest(os.Stdout, "authorize", c.Request)
		}
		// store, err := session.Start(c, c.Writer, c.Request)
		// if err != nil {
		// 	c.String(http.StatusInternalServerError, err.Error())
		// 	return
		// }
		// var form url.Values
		// if v, ok := store.Get("ReturnUri"); ok {
		// 	form = v.(url.Values)
		// }
		// c.Request.Form = form
		// store.Delete("ReturnUri")
		// store.Save()

		oauth.HandleAuthorizeRequest(c)

	})

	auth := g.Group("/oauth2")
	{
		// 应用程序通过此请求获取token
		auth.GET("/token", oauth.HandleTokenRequest)
	}

	api := g.Group("/api")
	{
		api.Use(oauth.HandleTokenRequest)
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

type LoginForm struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}

// oauth框架通过本方法识别用户身份信息,并且可以人为进行登录状态校验
// 本方法正常执行后,则会为客户端分配授权码(authorization_code)
func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "userAuthorizeHandler", r) // Ignore the error
	}
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

func authHandler(c *gin.Context) {
	if dumpvar {
		_ = dumpRequest(os.Stdout, "auth", c.Request) // Ignore the error
	}
	store, err := session.Start(nil, c.Writer, c.Request)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		c.Writer.Header().Set("Location", "/login")
		c.Writer.WriteHeader(http.StatusFound)
		return
	}
	c.HTML(http.StatusOK, "auth.html", gin.H{
		"title": "auth",
	})

}
