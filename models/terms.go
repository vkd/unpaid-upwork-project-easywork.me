package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type TermsSet struct {
	ID              bson.ObjectId `json:"id" db:"id"`
	CreatedDateTime time.Time     `json:"created_date_time" db:"created_date_time"`

	CurrencyId         string `json:"currency_id" db:"currency_id"`
	PricePerHour       int    `json:"price_per_hour" db:"price_per_hour"`
	WeeklyLimitInHours int    `json:"weekly_limit_in_hours" db:"weekly_limit_in_hours"`
}
