package httpservice

import (
	"log"

	"github.com/gin-gonic/gin"
	"gitlab.com/easywork.me/backend/models"
	"gitlab.com/easywork.me/backend/storage"
)

type Config struct {
	Addr      string
	SecretJWT []byte
}

// Start service
func Start(cfg Config, isDebug bool, db *storage.Storage) error {
	var r *gin.Engine
	if isDebug {
		r = gin.Default()
		r.Use(corsMiddleware())
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	// r.POST("/login", loginHandler(db))

	auth := r.Group("/", authMiddleware(cfg.SecretJWT))
	invitation := auth.Group("/invitations")
	{
		invitation.POST("/:id/accept", AccessRole(models.Work), invitationAcceptHandler(db))
		invitation.POST("/", AccessRole(models.Hire), invitationCreateHandler(db))
	}

	log.Printf("Web server is running on %s", cfg.Addr)
	return r.Run(cfg.Addr)
}

func apiError(c *gin.Context, code int, err interface{}) {
	c.Abort()
	c.JSON(code, err)
}
