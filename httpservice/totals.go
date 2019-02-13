package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gitlab.com/easywork.me/backend/storage"
)

func totalDailyHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		cID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		_, err = db.ContractGet(c, cID, user)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		from := c.Query("from")
		to := c.Query("to")

		out, err := db.TotalDaily(c, cID, from, to)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, out)
	}
}
