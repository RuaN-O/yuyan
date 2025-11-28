package auth

// Whitelist 免认证路由白名单
var Whitelist = map[string]bool{
	"/user.v1.UserService/Register": true,
	"/user.v1.UserService/Login":    true,
}

// IsWhitelist 检查路由是否在白名单中
func IsWhitelist(method string) bool {
	return Whitelist[method]
}
