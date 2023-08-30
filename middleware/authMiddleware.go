package middleware

/*
AuthMiddleware中间件
在每个需要验证登录状态的请求前进行 token 的验证.
如果 token 有效，就可以获取用户信息，并将用户信息添加到请求的上下文（context）中，方便后续处理。



*/
import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"time"
)

var Key = []byte("Reborn_But_In_Go") //加密key
type MyClaims struct {
	UserId   int64  `json:"user_id"`
	UserName string `json:"username"`
	jwt.StandardClaims
}

// CreateToken 生成一个token
func CreateToken(userId int64, userName string) (string, error) {
	expireTime := time.Now().Add(24 * time.Hour) //过期时间
	nowTime := time.Now()                        //当前时间
	claims := MyClaims{
		UserId:   userId,
		UserName: userName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), //过期时间戳
			IssuedAt:  nowTime.Unix(),    //当前时间戳
			Issuer:    "zhoumo",          //颁发者签名
			Subject:   "userToken",       //签名主题
		},
	}
	tokenStruct := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenStruct.SignedString(Key)
}

// CheckToken 验证token
func CheckToken(token string) (*MyClaims, bool) {
	tokenObj, _ := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Key, nil
	})
	if key, _ := tokenObj.Claims.(*MyClaims); tokenObj.Valid {
		return key, true
	} else {
		return nil, false
	}
}

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
