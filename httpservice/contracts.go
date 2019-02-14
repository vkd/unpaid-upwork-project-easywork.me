package httpservice

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"gitlab.com/easywork.me/backend/models"
	"gitlab.com/easywork.me/backend/storage"
)

func contractsGetHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		cs, err := db.ContractsGet(c, user)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, cs)
	}
}

func contractGetHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		cID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		contract, err := db.ContractGet(c, cID, user)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, contract)
	}
}

func contractCreateHandler(db *storage.Storage) gin.HandlerFunc {
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

func contractEndHandler(db *storage.Storage) gin.HandlerFunc {
	return contractChangeStatusHandler(db, models.Ended)
}
func contractPauseHandler(db *storage.Storage) gin.HandlerFunc {
	return contractChangeStatusHandler(db, models.Paused)
}
func contractResumeHandler(db *storage.Storage) gin.HandlerFunc {
	return contractChangeStatusHandler(db, models.Started, ContractResumed)
}

func contractChangeStatusHandler(db *storage.Storage, status models.ContractStatus, outs ...interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		cID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		err = db.ContractsUpdateStatus(c, cID, status, user)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		var out interface{}
		switch status {
		case models.Ended:
			out = ContractEnded
		case models.Paused:
			out = ContractPaused
		case models.Started:
			out = ContractStarted
		default:
			out = gin.H{"message": "status changed"}
		}

		if len(outs) > 0 {
			out = outs[0]
		}

		c.JSON(http.StatusOK, out)
	}
}
