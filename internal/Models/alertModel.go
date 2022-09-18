package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alert struct {
	ID     primitive.ObjectID `bson:"_id"`
	Url    string             `json:"url"`
	Name   string             `json:"name"`
	Time   string             `json:"time"`
	UserId string             `json:"user_id"`
}
