package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/easywork.me/backend/models"
	"gitlab.com/easywork.me/backend/storage"
)

func invitationCreateHandler(db *storage.Storage) gin.HandlerFunc {
	type InvitationCreate struct {
		models.InvitationBase
		models.TermsSetBase
	}
	return func(c *gin.Context) {
		user := getUser(c)

		var j InvitationCreate
		err := c.ShouldBindJSON(&j)
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		inv, err := db.InvitationCreate(c, &j.InvitationBase, &j.TermsSetBase, user.ID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, inv)
	}
}

func invitationAcceptHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)
		invID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		invitation, err := db.InvitationUpdateStatus(c, invID, user.ID, models.InvitationStatusAccepted)
		if err != nil {
			apiError(c, http.StatusUnprocessableEntity, err)
			return
		}

		cb := models.NewContractBase().FromInvitation(invitation)
		cb.Status = models.Started

		contract, err := db.ContractsCreate(c, cb, user.ID)
		if err != nil {
			apiError(c, http.StatusUnprocessableEntity, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"contract_id": contract.ID,
		})
	}
}
