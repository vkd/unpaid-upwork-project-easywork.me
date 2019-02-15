package httpservice

import (
	"log"
	"net/http"

	"github.com/gin-contrib/static"
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

	r.Use(static.Serve("/", static.LocalFile("./frontend/", false)))
	r.NoRoute(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Status(200)
		} else {
			c.File("./frontend/index.html")
		}
	})

	claimer := &claimCreator{secret: cfg.SecretJWT}

	r.GET("/status", statusHandler(db))
	r.POST("/users", userCreateHandler(db, claimer))
	r.POST("/user/login", userLoginHandler(db, claimer))

	authMid := authMiddleware(cfg.SecretJWT)

	auth := r.Group("", authMid)
	auth.PATCH("/profile", profileUpdateHandler(db))

	user := auth.Group("/user")
	{
		user.GET("", userGetHandler(db))
		user.PATCH("/password", changePasswordHandler(db))
		user.DELETE("", userDeleteHandler(db))
	}
	users := auth.Group("/users")
	{
		users.GET("", usersGetHandler(db))
		users.GET("/:user", userGetHandler(db))
	}
	tokens := auth.Group("/tokens")
	{
		tokens.POST("/verify", tokensVerifyHandler(db, cfg.SecretJWT))
	}
	invitations := auth.Group("/invitations")
	{
		invitations.GET("", invitationsGetHandler(db))
		invitations.POST("", AccessRole(models.Hire), invitationCreateHandler(db))
		invitationID := invitations.Group("/:id")
		{
			invitationID.GET("", invitationGetHandler(db))
			invitationID.POST("/accept", AccessRole(models.Work), invitationAcceptHandler(db))
			invitationID.POST("/decline", AccessRole(models.Work), invitationDeclineHandler(db))
			invitationID.DELETE("", AccessRole(models.Hire), invitationDeleteHandler(db))
		}
	}
	projects := auth.Group("/projects")
	{
		projects.GET("", projectsGetHandler(db))
		projects.POST("", AccessRole(models.Hire), projectCreateHandler(db))
		projectID := projects.Group("/:id")
		{
			projectID.GET("", projectGetHandler(db))
			projectID.DELETE("", AccessRole(models.Hire), projectDeleteHandler(db))
		}
	}
	contracts := auth.Group("/contracts")
	{
		contracts.GET("", contractsGetHandler(db))
		contracts.POST("", AccessRole(models.Hire), contractCreateHandler(db))

		contractID := contracts.Group("/:id")
		{
			contractID.GET("", contractGetHandler(db))
			contractID.GET("/dailies", totalDailyHandler(db))
			contractID.GET("/totals", totalsGetHandler(db))

			contractID.POST("/end", contractEndHandler(db))
			contractID.POST("/pause", contractPauseHandler(db))
			contractID.POST("/resume", contractResumeHandler(db))

			events := contractID.Group("/events")
			{
				events.GET("", eventsGetHandler(db))
				events.POST("/:type", AccessRole(models.Work), eventCreateHandler(db))
			}
		}
	}

	log.Printf("Web server is running on %s", cfg.Addr)
	return r.Run(cfg.Addr)
}

func statusHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	}
}

func apiError(c *gin.Context, code int, obj interface{}) {
	c.Abort()

	if obj == nil {
		c.JSON(code, gin.H{"error": "unknown error"})
		return
	}

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
