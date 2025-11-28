package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"

	// 引入你的配置包
	"yuyan/internal/conf"
)

// JwtRepo 是 JWT 生成和解析的接口，方便后续测试和扩展
type JwtRepo interface {
	GenerateToken(userID int64, username string) (string, error)
	ParseToken(tokenstring string) (*jwt.MapClaims, error)
}

// jwtData 是 JwtRepo 的实现，它持有配置
type jwtData struct {
	secret  string
	expire  time.Duration
	refresh time.Duration
}

// NewJwtRepo 是构造函数，用于创建 jwtData 实例
func NewJwtRepo(c *conf.Jwt) JwtRepo {
	// 注意：这里需要将 Protobuf 的 Duration 转换为 Go 的 time.Duration
	return &jwtData{
		secret:  c.GetSecret(),
		expire:  c.GetExpire().AsDuration(),
		refresh: c.GetRefresh().AsDuration(),
	}
}

// GenerateToken 现在是 jwtData 的一个方法，可以直接使用自己持有的配置
func (j *jwtData) GenerateToken(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		// 使用从配置中加载的过期时间
		"exp": time.Now().Add(j.expire).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用从配置中加载的密钥
	return token.SignedString([]byte(j.secret))
}

// 你可以在这里添加 ParseToken 等其他方法
func (j *jwtData) ParseToken(tokenString string) (*jwt.MapClaims, error) {
	// 使用 jwt.ParseWithClaims 来解析 token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		// 返回密钥
		return []byte(j.secret), nil
	})

	if err != nil {
		// 如果解析失败（如过期、签名无效），返回错误
		return nil, err
	}

	// 检查 token 是否有效，并断言 claims
	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
