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

func eventsGetHandler(db *storage.Storage) gin.HandlerFunc {
	return func(c *gin.Context) {
		cID, err := ObjectIDParam(c, "id")
		if err != nil {
			apiError(c, http.StatusBadRequest, err)
			return
		}

		var from *time.Time
		if s := c.Query("from"); s != "" {
			t, err := time.Parse("2006-01-02", s)
			if err != nil {
				apiError(c, http.StatusBadRequest, err)
				return
			}
			from = &t
		}

		var to *time.Time
		if s := c.Query("to"); s != "" {
			t, err := time.Parse("2006-01-02", s)
			if err != nil {
				apiError(c, http.StatusBadRequest, err)
				return
			}
			to = &t
		}

		out, err := db.EventsGet(c, cID, nil, from, to)
		if err != nil {
			apiError(c, http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, out)
	}
}

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
			e.ScreenshotFilename = time.Now().Format("150405") + ".jpg"
		}
		e.CreatedDateTime = time.Now()
		dayFolder := e.CreatedDateTime.Format("2006-01-02")
		e.ScreenshotUrl = fmt.Sprintf(`https://s3-us-west-2.amazonaws.com/hourly-tracker/%s/%s/%s/%s`, user.ID, contractID.Hex(), dayFolder, e.ScreenshotFilename)

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

		if eventType == models.EventLog {
			err = db.TotalsUpdate(c, contractID, event.CreatedDateTime, ValuePerEvent)
			if err != nil {
				apiError(c, http.StatusInternalServerError, err)
				return
			}
		}

		c.JSON(http.StatusOK, event)
	}
}
