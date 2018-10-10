package main

import (
	"gopkg.in/mgo.v2/bson"
)

// Item - Our struct for all items
type Item struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	Value       string        `bson:"value" json:"value"`
	Description string        `bson:"description" json:"description"`
}

type Config struct {
	Hosts      string `json:"hosts"`
	Database   string `json:"database"`
	Collection string `json:"collection"`
}
