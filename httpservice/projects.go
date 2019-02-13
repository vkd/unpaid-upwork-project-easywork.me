package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/easywork.me/backend/models"
	"gitlab.com/easywork.me/backend/storage"
)

func projectCreateHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		var j models.ProjectBase
		if err := c.ShouldBindJSON(&j); err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		p, err := db.ProjectCreate(c, &j, user.ID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, p)
	}
}

func projectDeleteHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		pID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		err = db.ProjectDelete(c, pID, user.ID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, ProjectDeleted)
	}
}
