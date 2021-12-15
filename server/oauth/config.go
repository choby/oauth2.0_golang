package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
)

func SetTokenType(tokenType string) {
	oauth2Server.Config.TokenType = tokenType
}

func SetAllowGetAccessRequest(allow bool) {
	oauth2Server.Config.AllowGetAccessRequest = allow
}

func SetAllowedResponseType(types ...oauth2.ResponseType) {
	oauth2Server.Config.AllowedResponseTypes = types
}

func SetAllowedGrantType(types ...oauth2.GrantType) {
	oauth2Server.Config.AllowedGrantTypes = types
}

func SetClientInfoHandler(handler server.ClientInfoHandler) {
	oauth2Server.ClientInfoHandler = handler
}

func SetClientAuthorizedHandler(handler server.ClientAuthorizedHandler) {
	oauth2Server.ClientAuthorizedHandler = handler
}

func SetClientScopeHandler(handler server.ClientScopeHandler) {
	oauth2Server.ClientScopeHandler = handler
}

func SetUserAuthorizationHandler(handler server.UserAuthorizationHandler) {
	oauth2Server.UserAuthorizationHandler = handler
}

func SetPasswordAuthorizationHandler(handler server.PasswordAuthorizationHandler) {
	oauth2Server.PasswordAuthorizationHandler = handler
}

func SetRefreshingScopeHandler(handler server.RefreshingScopeHandler) {
	oauth2Server.RefreshingScopeHandler = handler
}

func SetInternalErrorHandler(handler server.InternalErrorHandler) {
	oauth2Server.InternalErrorHandler = handler
}

func SetExtensionFieldsHandler(handler server.ExtensionFieldsHandler) {
	oauth2Server.ExtensionFieldsHandler = handler
}

func SetAccessTokenExpHandler(handler server.AccessTokenExpHandler) {
	oauth2Server.AccessTokenExpHandler = handler
}

func SetAuthorizeScopeHandler(handler server.AuthorizeScopeHandler) {
	oauth2Server.AuthorizeScopeHandler = handler
}
