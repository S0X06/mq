package middleware

import (
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

const (
	secret = "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq2iBb5"
)

var jwtSecret = []byte(secret)

type Claims struct {
	Data interface{} `json:"data"`
	jwt.StandardClaims
}

/**
* 生成token
 */
func GenerateToken(data interface{}) (token string, err error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(2 * time.Hour)

	claims := Claims{
		data,
		jwt.StandardClaims{
			Audience:  "user",            // 受众
			ExpiresAt: expireTime.Unix(), // 失效时间
			// Id:        string(user.ID),   // 编号
			IssuedAt:  nowTime.Unix(), // 签发时间
			Issuer:    "gin hello",    // 签发人
			NotBefore: nowTime.Unix(), // 生效时间
			Subject:   "login",        // 主题
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenClaims.SignedString(jwtSecret)
	token = "Bearer " + token
	return
}

// 解码
func ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}

//gin 中间键
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {

		var code int = 0
		var message string = "ok"

		token := c.Request.Header.Get("token")
		if token == "" {
			code = 10000
			message = "token 不能为空"
			c.JSON(http.StatusUnauthorized, gin.H{"code": code, "message": message})
			c.Abort()
			return
		}
		claims, err := ParseToken(token)
		if err != nil {
			code = 10001
			message = "token 无效"
			c.JSON(http.StatusUnauthorized, gin.H{"code": code, "message": message})
			c.Abort()
			return
		}
		if time.Now().Unix() > claims.ExpiresAt {
			code = 10002
			message = "token 超时"
			c.JSON(http.StatusUnauthorized, gin.H{"code": code, "message": message})
			c.Abort()
			return
		}

		c.Next()
	}
}

//解码用户信息
func ParseData(c *gin.Context) (data map[string]interface{}, err error) {
	token := c.Request.Header.Get("token")
	claims := &Claims{}
	claims, err = ParseToken(token)
	if err != nil {
		return
	}

	data = claims.Data.(map[string]interface{})
	return
}
