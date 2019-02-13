package httpservice

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gitlab.com/easywork.me/backend/filestorage"
	"gitlab.com/easywork.me/backend/models"
	"gitlab.com/easywork.me/backend/storage"
)

const (
	ValuePerEvent = 10
)

func eventCreateHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		user := getUser(c)

		contractID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		_, err = db.ContractGet(c, contractID, user)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		eventType := models.EventType(c.Param("type"))
		if err = models.CheckEventType(eventType); err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		e := models.NewEventBase()
		e.ContractID = contractID
		e.EventType = eventType
		e.Title = c.Request.FormValue("title")
		if e.ScreenshotFilename == "" {
			e.ScreenshotFilename = strconv.FormatInt(time.Now().Unix(), 10) + "HH24MISS.jpg"
		}
		e.CreatedDateTime = time.Now()

		if eventType == models.EventLog {
			e.KeyboardEventsCount, err = strconv.Atoi(c.Request.FormValue("keyboard_events_count"))
			if err != nil {
				apiError(c, http.StatusBadRequest, errors.Wrapf(err, "'keyboard_events_count' is not a number"))
				return
			}
			e.MouseEventsCount, err = strconv.Atoi(c.Request.FormValue("mouse_events_count"))
			if err != nil {
				apiError(c, http.StatusBadRequest, errors.Wrapf(err, "'mouse_events_count' is not a number"))
				return
			}

			dayFolder := e.CreatedDateTime.Format("2006-01-02")
			key := fmt.Sprintf(`%s/%s/%s/%s`, user.ID, contractID.Hex(), dayFolder, e.ScreenshotFilename)
			e.ScreenshotUrl = key
		}

		event, err := db.EventCreate(c, e, user)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		if eventType == models.EventLog {
			file, header, err := c.Request.FormFile("screenshot_file")
			if err != nil && err != http.ErrMissingFile {
				apiError(c, http.StatusBadRequest, errors.Wrapf(err, "is error on load 'screenshot_file'"))
				return
			}
			if err == nil {
				defer file.Close()

				err = filestorage.UploadFileToS3(file, header.Size, "hourly-tracker", e.ScreenshotUrl)
				if err != nil {
					apiError(c, http.StatusInternalServerError, errors.Wrapf(err, "error on upload to aws s3"))
					return
				}
			}
		}

		err = db.TotalsUpdate(c, contractID, event.CreatedDateTime, ValuePerEvent)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, event)
	}
}
