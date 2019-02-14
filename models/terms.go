package models

import (
	"github.com/pkg/errors"
)

// type TermsSet struct {
// 	ID primitive.ObjectID `json:"id" bson:"id"`

// 	TermsBase
// }

type TermsBase struct {
	// CreatedDateTime time.Time `json:"created_date_time" bson:"created_date_time"`
	// InvitationID    primitive.ObjectID `json:"invitation_id" bson:"invitation_id"`

	Currency           Currency `json:"currency_id" bson:"currency_id"`
	PricePerHour       uint     `json:"price_per_hour" bson:"price_per_hour"`
	WeeklyLimitInHours uint     `json:"weekly_limit_in_hours" bson:"weekly_limit_in_hours"`
}

func (t *TermsBase) PreCreate() *TermsBase {
	if t.Currency == "" {
		t.Currency = EUR
	}
	return t
}

func (t *TermsBase) Validate() error {
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
