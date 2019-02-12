package httpservice

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
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

	claimer := &claimCreator{secret: cfg.SecretJWT}

	// CreateUser
	r.POST("/users", userCreateHandler(db, claimer)) // +

	auth := r.Group("/", authMiddleware(cfg.SecretJWT))
	invitations := auth.Group("/invitations")
	{
		invitations.POST("/:id/accept", AccessRole(models.Work), invitationAcceptHandler(db))
		invitations.POST("/", AccessRole(models.Hire), invitationCreateHandler(db))
	}
	projects := auth.Group("/projects")
	{
		projects.POST("/", projectCreateHandler(db))
	}

	log.Printf("Web server is running on %s", cfg.Addr)
	return r.Run(cfg.Addr)
}

func apiError(c *gin.Context, code int, obj interface{}) {
	c.Abort()

	err, ok := obj.(error)
	if !ok {
		c.JSON(code, obj)
		return
	}

	switch e := errors.Cause(err).(type) {
	case *models.TrackerError:
		c.JSON(code, e)
	default:
		c.JSON(code, gin.H{
			"error": err.Error(),
		})
	}
}
