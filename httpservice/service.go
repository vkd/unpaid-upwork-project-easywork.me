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

	user := r.Group("/user")
	{
		// DeleteUser
		user.DELETE("/user", userDeleteHandler(db))
	}
	invitations := auth.Group("/invitations")
	{
		// CreateInvitation
		invitations.POST("/", AccessRole(models.Hire), invitationCreateHandler(db))
		invitationID := invitations.Group("/:id")
		{
			// AcceptInvitation
			invitationID.POST("/accept", AccessRole(models.Work), invitationAcceptHandler(db))
			// DeclineInvitation
			invitationID.POST("/decline", AccessRole(models.Work), invitationDeclineHandler(db))
			// DeleteInvitation
			invitationID.DELETE("/", AccessRole(models.Hire), invitationDeleteHandler(db))
		}
	}
	projects := auth.Group("/projects")
	{
		// CreateProject
		projects.POST("/", AccessRole(models.Hire), projectCreateHandler(db))
		projectID := projects.Group("/:id")
		{
			// DeleteProject
			projectID.DELETE("/", AccessRole(models.Hire), projectDeleteHandler(db))
		}
	}
	contracts := auth.Group("/contracts")
	{
		// CreateContract
		contracts.POST("/", AccessRole(models.Hire), contractCreateHandler(db))

		contractID := contracts.Group("/:id")
		{
			contractID.GET("/dailies", totalDailyHandler(db))
			// EndContract
			contractID.POST("/end", AccessRole(models.Hire), contractEndHandler(db))
			events := contractID.Group("/events")
			{
				// CreateContractEvent
				events.POST("/:type", AccessRole(models.Work), eventCreateHandler(db))
			}
		}
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
	case storage.ErrNoUpdated:
		code = http.StatusNotModified
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
