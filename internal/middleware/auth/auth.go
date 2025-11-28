package auth

import (
	"context"
	"strconv"

	// 使用别名避免和标准库 errors 冲突
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"

	//"github.com/golang-jwt/jwt/v4"

	"yuyan/internal/pkg/jwt"
)

// contextKey 是一个自定义类型，用于防止 context 中的 key 冲突
type contextKey string

const (
	// UserIDKey 是在 context 中存储用户ID的key
	UserIDKey contextKey = "user_id"
	// AuthorizationKey 认证头键名
	AuthorizationKey = "Authorization"
	// BearerWord Bearer 标识
	BearerWord = "Bearer"
)

var (
	ErrMissingToken = errors.Unauthorized("UNAUTHORIZED", "missing authorization token")
	ErrInvalidToken = errors.Unauthorized("UNAUTHORIZED", "invalid authorization token")
)

type AuthMiddleware struct {
	jwt jwt.JwtRepo
}

func NewAuthMiddleware(jwtrepo jwt.JwtRepo) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwtrepo,
	}
}

// NewAuthFunc 返回一个符合 Kratos auth 中间件要求的验证函数
// 它接收一个 JwtRepo 实例，用于实际的 token 解析
func (m *AuthMiddleware) Handler() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 检查是否在白名单中
			if tr, ok := transport.FromServerContext(ctx); ok {
				log.Infof("看看路径:%v", tr.Operation())
				if IsWhitelist(tr.Operation()) {
					return handler(ctx, req)
				}

				authHeader := tr.RequestHeader().Get(AuthorizationKey)
				if authHeader == "" {
					return nil, ErrMissingToken
				}

				tokenString, ok := parseBearerToken(authHeader)
				if !ok {
					return nil, ErrInvalidToken
				}

				claims, err := m.jwt.ParseToken(tokenString)
				log.Infof("看看claims:%v\n", claims)
				if err != nil || claims.Valid() != nil {
					return nil, ErrInvalidToken
				}

				// 将用户信息存入上下文
				var userID int64
				userid, ok := (*claims)["user_id"]
				if !ok {
					log.Errorf("用户id解析失败:%v\n", userid)
					return nil, ErrInvalidToken
				} else {
					userID, ok := userid.(float64)
					if !ok {
						log.Errorf("用户id解析失败step2:%v\n", userID)
					}
				}
				username := (*claims)["username"].(string)
				log.Infof("看看username:%v", username)
				ctx = context.WithValue(ctx, UserIDKey, userID)
				ctx = context.WithValue(ctx, contextKey("username"), username)
				//将用户信息存入元信息
				if md, ok := metadata.FromServerContext(ctx); ok {
					md.Set("yuyan-user-id", strconv.Itoa(int(userID)))
					md.Set("yuyan-user-name", username)
				}
			}

			return handler(ctx, req)
		}
	}
}

// FromContext 从 context 中获取 user_id
func FromContext(ctx context.Context) (int64, bool) {
	v := ctx.Value(UserIDKey)
	if userID, ok := v.(int64); ok {
		return userID, true
	}
	return 0, false
}

// MustFromContext 从 context 中获取 user_id，如果不存在则 panic
func MustFromContext(ctx context.Context) int64 {
	userID, ok := FromContext(ctx)
	if !ok {
		panic("user_id not found in context")
	}
	return userID
}
func parseBearerToken(authHeader string) (string, bool) {
	return authHeader, true
}
