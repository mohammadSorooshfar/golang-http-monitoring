package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Url struct {
	ID        primitive.ObjectID `bson:"_id"`
	Link      string             `json:"link"`
	Success   int                `json:"success"`
	Failed    int                `json:"failed"`
	User_id   string             `json:"user_id"`
	Threshold int                `json:"treshold"`
}
