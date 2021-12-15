package oauth

import (
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
)

func SetTokenType(tokenType string) {
	srv.Config.TokenType = tokenType
}

func SetAllowGetAccessRequest(allow bool) {
	srv.Config.AllowGetAccessRequest = allow
}

func SetAllowedResponseType(types ...oauth2.ResponseType) {
	srv.Config.AllowedResponseTypes = types
}

func SetAllowedGrantType(types ...oauth2.GrantType) {
	srv.Config.AllowedGrantTypes = types
}

func SetClientInfoHandler(handler server.ClientInfoHandler) {
	srv.ClientInfoHandler = handler
}

func SetClientAuthorizedHandler(handler server.ClientAuthorizedHandler) {
	srv.ClientAuthorizedHandler = handler
}

func SetClientScopeHandler(handler server.ClientScopeHandler) {
	srv.ClientScopeHandler = handler
}

func SetUserAuthorizationHandler(handler server.UserAuthorizationHandler) {
	srv.UserAuthorizationHandler = handler
}

func SetPasswordAuthorizationHandler(handler server.PasswordAuthorizationHandler) {
	srv.PasswordAuthorizationHandler = handler
}

func SetRefreshingScopeHandler(handler server.RefreshingScopeHandler) {
	srv.RefreshingScopeHandler = handler
}

func SetInternalErrorHandler(handler server.InternalErrorHandler) {
	srv.InternalErrorHandler = handler
}

func SetExtensionFieldsHandler(handler server.ExtensionFieldsHandler) {
	srv.ExtensionFieldsHandler = handler
}

func SetAccessTokenExpHandler(handler server.AccessTokenExpHandler) {
	srv.AccessTokenExpHandler = handler
}

func SetAuthorizeScopeHandler(handler server.AuthorizeScopeHandler) {
	srv.AuthorizeScopeHandler = handler
}
