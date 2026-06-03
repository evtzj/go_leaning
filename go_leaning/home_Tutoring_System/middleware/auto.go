package middleware

import (
	"net/http"
	"strings"

	"home_Tutoring_System/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Claims JWT 的 payload 结构体。前端登录后拿到的 token 里就存了这些字段。
// 嵌入 jwt.RegisteredClaims 等于继承了标准字段：ExpiresAt（过期时间）、IssuedAt（签发时间）等。
// Claims是一个通行证的结构体
type Claims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// AuthRequired 是一个 Gin 中间件，用来验证请求里的 JWT token。
// 返回 gin.HandlerFunc → 满足了 r.Use() 的参数要求。
//
// 使用方式（在 router 里）：
//
//	auth := r.Group("/api")
//	auth.Use(middleware.AuthRequired())  // 这个组里的所有接口都需要登录
func AuthRequired() gin.HandlerFunc {
	// 返回一个闭包函数，这才是真正处理请求的中间件
	return func(c *gin.Context) {

		// ── 第 1 步：从请求头里取 Authorization 字段 ──────────────────
		// HTTP 请求头长这样：Authorization: Bearer eyJhbGciOi...
		//这个函数相当于一个保安,现在需要向客人要通行证
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// 连 Authorization 头都没有 → 直接返回 401，不再往下走
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "请先登录",
			})
			return // ⚠️ 这个 return 很重要，防止继续执行后面的代码
		}

		// ── 第 2 步：解析 "Bearer <token>" 格式 ──────────────────────
		// SplitN 按空格切，最多切 2 段："Bearer" 和 "eyJhbG..."
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "认证格式错误",
			})
			return
		}
		tokenString := parts[1] // 第二段就是 token 本身

		// ── 第 3 步：用密钥验证 token 签名 ──────────────────────────
		claims := &Claims{} // 传指针进去，验签成功后会填上值
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			// 这个回调函数告诉 JWT 库"用什么密钥来验签"
			// 必须用和签发 token 时同样的密钥（config.JWTSecret）
			return config.JWTSecret, nil
		})

		// 验签失败 or token 过期 → 401
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"message": "Token 无效或已过期",
			})
			return
		}

		// ── 第 4 步：把用户信息注入 context，供下游 handler 使用 ────
		// 到了这里说明 token 合法，claims 里已经有值了
		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		// ── 第 5 步：放行 ────────────────────────────────────────────
		// 调用 c.Next() 的意思是"我这边没问题了，交给下一个中间件或 handler"
		c.Next()
	}
}
