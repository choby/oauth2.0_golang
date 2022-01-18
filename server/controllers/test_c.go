package controllers

import (
	"net/http"
	"time"

	"github.com/choby/oauth2.0_golang/server/oauth"
	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	token, err := oauth.SRV.ValidationBearerToken(c.Request) // 验证token, 可以用这个方法抽象一个action用于资源服务器调用验证token有效性
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// cli, err := mgr.GetClient(token.GetClientID())
	// if err != nil {
	//     http.Error(w, err.Error(), http.StatusBadRequest)
	//     return
	// }

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"user_id":    token.GetUserID(),
		"client_id":  token.GetClientID(),
		"scope":      token.GetScope(),
		// "domain": cli.GetDomain(),
	}
	c.JSON(http.StatusOK, data)
}
