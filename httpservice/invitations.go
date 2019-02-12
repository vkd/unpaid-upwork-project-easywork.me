package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"

	"gitlab.com/easywork.me/backend/models"
	"gitlab.com/easywork.me/backend/storage"
)

func invitationCreateHandler(db *storage.Storage) gin.HandlerFunc {
	type invitationCreate struct {
		InvitationID string `json:"invitation_id"`
		models.TermsSet
	}
	return func(c *gin.Context) {
		user := getUser(c)

		var j invitationCreate
		err := c.ShouldBindJSON(&j)
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		inv, err := db.InvitationCreate(c, bson.ObjectIdHex(j.InvitationID), user.ID, &j.TermsSet, &user)
		if err != nil {
			apiError(c, http.StatusUnprocessableEntity, err)
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

		contract := models.NewContract().FromInvitation(invitation)
		contract.Status = models.Started

		contract, err = db.ContractsCreate(c, contract)
		if err != nil {
			apiError(c, http.StatusUnprocessableEntity, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"contract_id": contract.ID,
		})
	}
}
