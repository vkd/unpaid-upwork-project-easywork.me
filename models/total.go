package models

import (
	"github.com/mongodb/mongo-go-driver/bson/primitive"
)

type Total struct {
	ContractID primitive.ObjectID `json:"contract_id" bson:"contract_id"`
	Date       string             `json:"date" bson:"date"`
	Value      int                `json:"value" bson:"value"`
}
