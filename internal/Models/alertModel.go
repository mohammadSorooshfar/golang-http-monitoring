package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Alert struct {
	ID     primitive.ObjectID `bson:"_id"`
	Url    string             `query:"url"on:"url"`
	Owner  string             `json:"owner"`
	Time   string             `json:"time"`
	UserId primitive.ObjectID `json:"user_id"`
}
