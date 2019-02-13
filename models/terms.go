package models

import (
	"time"

	"github.com/mongodb/mongo-go-driver/bson/primitive"
	"github.com/pkg/errors"
)

type TermsSet struct {
	ID primitive.ObjectID `json:"id" bson:"id"`

	TermsSetBase
}

type TermsSetBase struct {
	CreatedDateTime time.Time `json:"created_date_time" bson:"created_date_time"`

	Currency           Currency `json:"currency_id" bson:"currency_id"`
	PricePerHour       uint     `json:"price_per_hour" bson:"price_per_hour"`
	WeeklyLimitInHours uint     `json:"weekly_limit_in_hours" bson:"weekly_limit_in_hours"`
}

func (t *TermsSetBase) PreCreate() *TermsSetBase {
	if t.Currency == "" {
		t.Currency = EUR
	}
	return t
}

func (t *TermsSetBase) Validate() error {
	if t.PricePerHour <= 0 || t.PricePerHour >= 500 {
		return errors.Errorf("'price_per_hour' not valid (allow > 0 and < 500)")
	}
	return nil
}

type Currency string

const (
	EUR Currency = "eur"
	USD Currency = "usd"
)
