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

	r.POST("/users", userCreateHandler(db, claimer))

	auth := r.Group("/", authMiddleware(cfg.SecretJWT))

	user := auth.Group("/user")
	{
		user.GET("/", userGetHandler(db))
		user.DELETE("/", userDeleteHandler(db))
	}
	users := auth.Group("/users")
	{
		users.GET("/", usersGetHandler(db))
	}
	invitations := auth.Group("/invitations")
	{
		invitations.GET("/", invitationsGetHandler(db))
		invitations.POST("/", AccessRole(models.Hire), invitationCreateHandler(db))
		invitationID := invitations.Group("/:id")
		{
			invitationID.GET("/", invitationGetHandler(db))
			invitationID.POST("/accept", AccessRole(models.Work), invitationAcceptHandler(db))
			invitationID.POST("/decline", AccessRole(models.Work), invitationDeclineHandler(db))
			invitationID.DELETE("/", AccessRole(models.Hire), invitationDeleteHandler(db))
		}
	}
	projects := auth.Group("/projects")
	{
		projects.GET("/", projectsGetHandler(db))
		projects.POST("/", AccessRole(models.Hire), projectCreateHandler(db))
		projectID := projects.Group("/:id")
		{
			projectID.GET("/", projectGetHandler(db))
			projectID.DELETE("/", AccessRole(models.Hire), projectDeleteHandler(db))
		}
	}
	contracts := auth.Group("/contracts")
	{
		contracts.GET("/", contractsGetHandler(db))
		contracts.POST("/", AccessRole(models.Hire), contractCreateHandler(db))

		contractID := contracts.Group("/:id")
		{
			contractID.GET("/", contractGetHandler(db))
			contractID.GET("/dailies", totalDailyHandler(db))
			contractID.GET("/totals", totalsGetHandler(db))
			contractID.POST("/end", AccessRole(models.Hire), contractEndHandler(db))

			events := contractID.Group("/events")
			{
				events.GET("/", eventsGetHandler(db))
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
