package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/heroku/go-base/models"
)

var identityKey = "ID"
var role = "role"

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// GetJWTMiddleware return a middleware with utility function for routing
func GetJWTMiddleware(RedisConnection redis.Conn) (*jwt.GinJWTMiddleware, error) {
	authMiddleware, _ := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
					role:        v.Grant,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			var user models.User

			s, err := redis.String(RedisConnection.Do("GET", uint(claims[identityKey].(float64))))

			if err != nil {
				user, _ = models.GetUserByID(uint(claims[identityKey].(float64)))
				fmt.Println("User does not exist")
			} else {
				err = json.Unmarshal([]byte(s), &user)
				fmt.Println(s)
			}

			return user
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			if customer, err := models.Authenticate(loginVals.Username, loginVals.Password); err != nil {
				return nil, jwt.ErrFailedAuthentication
			} else {
				basicUser := customer.GetBasicUser()
				json, _ := json.Marshal(basicUser)
				RedisConnection.Do("SET", basicUser.ID, json)
				// Here is a good place to push the customer into redis or similar
				return &customer, nil
			}
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if _, ok := data.(models.User); ok { //&& v.UserName == "admin" {
				return true
			}

			return false

		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		// - "param:<name>"
		TokenLookup: "header: Authorization, query: token, cookie: jwt",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	})
	return authMiddleware, nil
}

func GetAdminHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		currentRole := models.Grant(claims[role].(float64))
		if currentRole != models.Admin {
			c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "error": "You must be admin to use this resource"})
			c.Abort()
			return
		}
		c.Next()
	}
}
