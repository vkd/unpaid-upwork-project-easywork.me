package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/easywork.me/backend/models"
)

func authMiddleware(secretJWT []byte) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			apiError(c, http.StatusUnauthorized, &models.AccessForbidden)
			return
		}

		user, err := VerifyTokenString(authHeader, secretJWT)
		if err != nil {
			apiError(c, http.StatusUnauthorized, err)
			return
		}
		if user == nil {
			apiError(c, http.StatusUnauthorized, &models.AccessForbidden)
			return
		}

		setUser(c, user)
		c.Next()
	}
}

func getUser(c *gin.Context) (out *models.User) {
	u, ok := c.Get("user")
	if !ok {
		return &models.User{}
	}
	user, ok := u.(*models.User)
	if !ok {
		return &models.User{}
	}
	return user
}

func setUser(c *gin.Context, user *models.User) {
	if user == nil {
		panic("user is nil") // impossible way
	}
	c.Set("user", user)
}

func AccessRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)
		for _, r := range roles {
			if r == user.Role {
				c.Next()
				return
			}
		}
		apiError(c, http.StatusUnauthorized, &models.AccessForbidden)
	}
}
