package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/pkg/errors"
)

type EventType string

const (
	EventStart EventType = "start"
	EventLog   EventType = "log"
	EventStop  EventType = "stop"
)

func CheckEventType(et EventType) error {
	switch et {
	case EventStart, EventStop, EventLog:
		return nil
	}
	return errors.Errorf("wrong event type (allowed: start|stop|log)")
}

type Event struct {
	ID primitive.ObjectID `json:"id" bson:"id"`

	EventBase `bson:"inline"`
}

type EventBase struct {
	CreatedDateTime time.Time          `json:"created_date_time" bson:"created_date_time"`
	ContractID      primitive.ObjectID `json:"contract_id" bson:"contract_id"`
	EventType       EventType          `json:"event_type" bson:"event_type"`

	KeyboardEventsCount int    `json:"keyboard_events_count" bson:"keyboard_events_count"`
	MouseEventsCount    int    `json:"mouse_events_count" bson:"mouse_events_count"`
	ScreenshotFilename  string `json:"-" bson:"screenshot_filename"`
	ScreenshotUrl       string `json:"screenshot_url" bson:"screenshot_url"`
	Title               string `json:"title" bson:"title"`
}

func NewEventBase() *EventBase {
	return &EventBase{}
}
