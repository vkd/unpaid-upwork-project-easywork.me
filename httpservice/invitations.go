package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"gitlab.com/easywork.me/backend/models"
	"gitlab.com/easywork.me/backend/storage"
)

func invitationsGetHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		invs, err := db.InvitationsGet(c, user.ID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, invs)
	}
}

func invitationGetHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		iID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		inv, err := db.InvitationGet(c, iID, user.ID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, inv)
	}
}

func invitationCreateHandler(db *storage.Storage) gin.HandlerFunc {
	type InvitationCreate struct {
		models.InvitationBase
	}
	return func(c *gin.Context) {
		user := getUser(c)

		var j InvitationCreate
		err := c.ShouldBindJSON(&j)
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		if err = j.TermsBase.Validate(); err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		inv, err := db.InvitationCreate(c, &j.InvitationBase, user.ID)
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
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		cb := models.NewContractBase().FromInvitation(invitation)
		cb.Status = models.Started

		contract, err := db.ContractsCreate(c, cb, invitation.OwnerID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"contract_id": contract.ID,
		})
	}
}

func invitationDeclineHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		invID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		_, err = db.InvitationUpdateStatus(c, invID, user.ID, models.InvitationStatusDeclined)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, InvitationDeclined)
	}
}

func invitationDeleteHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		invID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		err = db.InvitationDelete(c, invID, user.ID)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, InvitationDeleted)
	}
}
