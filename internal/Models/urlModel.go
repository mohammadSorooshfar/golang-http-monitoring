package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Url struct {
	ID        primitive.ObjectID `bson:"_id"`
	Link      string             `json:"link"`
	Success   map[string]int     `json"success"`
	Failed    map[string]int
	User_id   string             `json:"user_id"`
	Threshold int                `json:"threshold"`
	Period    int                `json:"period"`
}
