package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func authMiddleware(c *gin.Context) {
	var userProfile *UserTokenWithClaims

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		var tErr *TrackerError
		userProfile, tErr = VerifyTokenString(authHeader)
		if tErr != nil {
			apiError(c, http.StatusUnauthorized, tErr)
			return
		}
	}

	var isAccessGranted bool
	if !isAccessGranted {
		apiError(c, http.StatusUnauthorized, &AccessForbidden)
		return
	}

	setUser(c, userProfile)
	c.Next()
}

func getUser(c *gin.Context) *UserTokenWithClaims {
	u, ok := c.Get("user")
	if !ok {
		return nil
	}
	user, ok := u.(*UserTokenWithClaims)
	if !ok {
		return nil
	}
	return user
}

func setUser(c *gin.Context, user *UserTokenWithClaims) {
	c.Set("user", user)
}
