package httpservice

import (
	"log"

	"github.com/gin-gonic/gin"
)

// Start service
func Start(addr string, isDebug bool) error {
	var r *gin.Engine
	if isDebug {
		r = gin.Default()
		r.Use(corsMiddleware())
	} else {
		r = gin.New()
	}

	r.Use(authMiddleware)

	// AcceptInvitation
	r.POST("/invitations/:id/accept", func(c *gin.Context) {
		c.String(200, "%s", "hello, world")
	})

	log.Printf("Web server is running on %s", addr)
	return r.Run(addr)
}

func apiError(c *gin.Context, code int, err interface{}) {
	c.Abort()
	c.JSON(code, err)
}
