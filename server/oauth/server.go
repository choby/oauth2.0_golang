package oauth

import (
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/choby/oauth2.0_golang/server/libs/request"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
	"github.com/golang-jwt/jwt"
)

var (
	SRV  *server.Server
	once sync.Once
)

func initClientStore() *store.ClientStore {
	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("222222", &models.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:9094",
	})

	return clientStore
}

func initManage() *manage.Manager {
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token存储方式
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// generate jwt access token
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	//manager.MapAccessGenerate(generates.NewAccessGenerate())

	clientStore := initClientStore()
	manager.MapClientStorage(clientStore)

	return manager
}

// 初始化服务
func InitServer() *server.Server {
	manager := initManage()
	once.Do(func() {
		// oauth2Server = server.NewServer(server.NewConfig(), manager)
		SRV = server.NewDefaultServer(manager)
		SRV.SetAllowedGrantType("authorization_code", "refresh_token")
		SRV.SetAllowGetAccessRequest(true)

		SRV.SetAuthorizeScopeHandler(authorizeScopeHandler)

		// 密码授权模式才需要用到这个配置, 这个模式不需要分配授权码,而是直接分配token,通常用于无后端的应用
		SRV.SetPasswordAuthorizationHandler(passwordAuthorizationHandler)

		// 这一行很关键,这个方法让oauth框架识别当前用户身份标识(并且可以人为处理登陆状态检验等等)
		// 具体看userAuthorizeHandler方法实现
		SRV.SetUserAuthorizationHandler(userAuthorizeHandler)

		SRV.SetInternalErrorHandler(internalErrorHandler)
		SRV.SetResponseErrorHandler(responseErrorHandler)
	})

	return SRV
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

func passwordAuthorizationHandler(username, password string) (userID string, err error) {
	if username == "test" && password == "test" {
		userID = "test"
	}
	return userID, nil
}

// 根据client注册的scope
// 过滤非法scope
func authorizeScopeHandler(w http.ResponseWriter, r *http.Request) (scope string, err error) {
	if r.Form == nil {
		r.ParseForm()
	}
	s := ScopeFilter(r.Form.Get("client_id"), r.Form.Get("scope"))
	if s == nil {
		http.Error(w, "Invalid Scope", http.StatusBadRequest)
		return
	}
	scope = ScopeJoin(s)

	return
}

func internalErrorHandler(err error) (re *errors.Response) {
	// log.Println("Internal Error:", err.Error())
	return
}

func responseErrorHandler(re *errors.Response) {
	// log.Println("Response Error:", re.Error.Error())
	return
}

type Scope struct {
	ID    string `yaml:"id"`
	Title string `yaml:"title"`
}

func ScopeJoin(scope []Scope) string {
	var s []string
	for _, sc := range scope {
		s = append(s, sc.ID)
	}
	return strings.Join(s, ",")
}

func ScopeFilter(clientID string, scope string) (s []Scope) {
	// cli := GetClient(clientID)
	// sl := strings.Split(scope, ",")
	// for _, str := range sl {
	// 	for _, sc := range cli.Scope {
	// 		if str == sc.ID {
	// 			s = append(s, sc)
	// 		}
	// 	}
	// }
	s = []Scope{}
	return
}
