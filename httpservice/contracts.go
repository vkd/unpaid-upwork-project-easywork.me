package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"gitlab.com/easywork.me/backend/models"
	"gitlab.com/easywork.me/backend/storage"
)

func contractsCreateHandler(db *storage.Storage) gin.HandlerFunc {
	type InvitationIdRequest struct {
		InvitationId primitive.ObjectID `json:"invitation_id"`
	}
	return func(c *gin.Context) {
		user := getUser(c)

		var j InvitationIdRequest
		if err := c.ShouldBindJSON(&j); err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		inv, err := db.InvitationGet(c, j.InvitationId, user.ID)
		if err != nil {
			if storage.IsNotFound(err) {
				apiError(c, http.StatusNotFound, &models.InvitationNotFound)
			} else {
				apiError(c, http.StatusInternalServerError, err)
			}
			return
		}

		cb := models.NewContractBase().FromInvitation(inv)
		cntr, err := db.ContractsCreate(c, cb, user.ID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, cntr)
	}
}
