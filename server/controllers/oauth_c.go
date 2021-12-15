package controllers

import (
	"net/http"
	"net/url"
	"os"

	"github.com/choby/oauth2.0_golang/server/libs/request"

	"github.com/choby/oauth2.0_golang/server/oauth"
	"github.com/gin-gonic/gin"
	"github.com/go-session/session"
)

var (
	dumpvar = true
)

type LoginForm struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Login(c *gin.Context) {
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

}

func Auth(c *gin.Context) {
	if dumpvar {
		_ = request.DumpRequest(os.Stdout, "auth", c.Request) // Ignore the error
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

func Authorize(c *gin.Context) {
	if dumpvar {
		request.DumpRequest(os.Stdout, "authorize", c.Request)
	}
	store, err := session.Start(c, c.Writer, c.Request)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
	}
	c.Request.Form = form
	store.Delete("ReturnUri")
	store.Save()

	oauth.HandleAuthorizeRequest(c)

}
