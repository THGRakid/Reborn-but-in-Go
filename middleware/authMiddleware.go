package middleware

/*
AuthMiddleware中间件
在每个需要验证登录状态的请求前进行 token 的验证.
如果 token 有效，就可以获取用户信息，并将用户信息添加到请求的上下文（context）中，方便后续处理。



*/
import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 提取 token
		token := c.Query("token")

		// 解析 token
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			// 验证签名密钥
			return []byte("Reborn_but_in_Go"), nil
		})
		if err == nil && parsedToken.Valid {
			// token 验证通过，继续处理请求
			c.Set("is_authenticated", true)
			return
		}

		// 验证失败，返回未授权错误
		c.Set("is_authenticated", false)
	}
}
