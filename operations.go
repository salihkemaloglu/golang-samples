package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Establish a connection to database
func Connect(connectionUrl string) {
	info := &mgo.DialInfo{
		Addrs:    []string{connectionUrl},
		Timeout:  5 * time.Second,
		Database: DB,
		Username: "",
		Password: "",
	}
	session, err := mgo.DialWithInfo(info)
	if err != nil {
		fmt.Println(err.Error())
	}
	db = session.DB(DB)
}

// Find list of Itemss
func FindAll() ([]Item, error) {
	var Items []Item
	err := db.C(COLLECTION).Find(bson.M{}).All(&Items)
	return Items, err
}

// Find a Items by its id
func FindById(id string) (Item, error) {
	var Items Item
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&Items)
	return Items, err
}

// Insert a Items into database
func Insert(Items Item) error {
	err := db.C(COLLECTION).Insert(&Items)
	return err
}

// Delete an existing Items
func Delete(Items Item) error {
	err := db.C(COLLECTION).Remove(&Items)
	return err
}

// Update an existing Items
func Update(Items Item) error {
	err := db.C(COLLECTION).Update(bson.M{"_id": Items.ID}, &Items)
	return err
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func LoadConfiguration() {
	var config Config
	configFile, err := os.Open("config.json")
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	var url = config.Hosts
	COLLECTION = config.Collection
	DB = config.Database
	Connect(url)
}
