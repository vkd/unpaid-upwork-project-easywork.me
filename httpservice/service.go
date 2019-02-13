package httpservice

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/mongo"
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
	r.POST("/users", userCreateHandler(db, claimer))

	auth := r.Group("/", authMiddleware(cfg.SecretJWT))
	invitations := auth.Group("/invitations")
	{
		// CreateInvitation
		invitations.POST("/", AccessRole(models.Hire), invitationCreateHandler(db))

		invitations.POST("/:id/accept", AccessRole(models.Work), invitationAcceptHandler(db))
	}
	projects := auth.Group("/projects")
	{
		// CreateProject
		projects.POST("/", AccessRole(models.Hire), projectCreateHandler(db))
	}
	contracts := auth.Group("/contracts")
	{
		// CreateContract
		contracts.POST("/", AccessRole(models.Hire), contractsCreateHandler(db))
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

	e := errors.Cause(err)

	switch e {
	case storage.ErrNotFound, mongo.ErrNoDocuments:
		code = http.StatusNotFound
	}

	switch e := e.(type) {
	case *models.TrackerError:
		c.JSON(code, e)
	default:
		c.JSON(code, gin.H{
			"error": err.Error(),
		})
	}
}
