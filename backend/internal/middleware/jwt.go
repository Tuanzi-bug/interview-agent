package middleware

import (
	"context"
	"errors"
	"strings"
	"time"

	"ai-eino-interview-agent/internal/config"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims 自定义JWT声明结构
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTMiddleware JWT认证中间件
func JWTMiddleware() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		// 获取Authorization头
		authHeader := string(ctx.GetHeader("Authorization"))
		if authHeader == "" {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Authorization header is required",
			})
			ctx.Abort()
			return
		}

		// 检查Bearer前缀
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Authorization header format must be Bearer {token}",
			})
			ctx.Abort()
			return
		}

		// 解析JWT token
		claims, err := parseToken(parts[1])
		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"code":    401,
				"message": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		// 将用户信息存储到上下文
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("role", claims.Role)

		ctx.Next(c)
	}
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username, role string) (string, error) {
	cfg := config.Global.Security

	// 解析过期时间
	expiration, err := time.ParseDuration(cfg.JWTExpiration)
	if err != nil {
		expiration = 24 * time.Hour // 默认24小时
	}

	// 创建声明
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "interview-agent",
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的编码后的字符串token
	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// parseToken 解析JWT token
func parseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(config.Global.Security.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// GetUserID 从上下文获取用户ID
func GetUserID(ctx *app.RequestContext) uint {
	userID, exists := ctx.Get("user_id")
	if !exists {
		return 0
	}
	return userID.(uint)
}

// GetUsername 从上下文获取用户名
func GetUsername(ctx *app.RequestContext) string {
	username, exists := ctx.Get("username")
	if !exists {
		return ""
	}
	return username.(string)
}

// GetUserRole 从上下文获取用户角色
func GetUserRole(ctx *app.RequestContext) string {
	role, exists := ctx.Get("role")
	if !exists {
		return ""
	}
	return role.(string)
}
