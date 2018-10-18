package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var db *mgo.Database
var ITEM string
var USER string
var DB string

// Find list of Items
func (r Item) FindAll() ([]Item, error) {
	var Items []Item
	err := db.C(ITEM).Find(bson.M{}).All(&Items)
	return Items, err
}

// Find a Items by its id
func (r Item) FindById() (Item, error) {
	err := db.C(ITEM).FindId(bson.ObjectIdHex(r.ItemId)).One(&r)
	return r, err
}

// Insert a Items into database
func (r Item) Insert() error {
	err := db.C(ITEM).Insert(&r)
	return err
}

// Delete an existing Items
func (r Item) Delete() error {
	err := db.C(ITEM).Remove(&r)
	return err
}

// Update an existing Items
func (r Item) Update() error {
	err := db.C(ITEM).Update(bson.M{"_id": r.ID}, &r)
	return err
}

// Find a user
func (r User) Login() (User, error) {
	err := db.C(USER).Find(bson.M{"username": r.Username, "password": r.Password}).One(&r)
	return r, err
}

func (r User) Register() bool {

	err := db.C("user").Insert(&r)
	if err != nil {
		return false
	}
	return true
}
func (r User) CheckUser() bool {
	err := db.C("user").Find(bson.M{"username": r.Username}).One(&r)
	if err != nil {
		return true
	} else {
		return false
	}
}
func (r User) FindAllUser() ([]User, error) {
	var users []User
	err := db.C(USER).Find(bson.M{}).All(&users)
	return users, err
}

// Delete an existing Items
func (r User) DeleteUser() error {
	err := db.C(USER).Remove(&r)
	return err
}

// Update an existing Items
func (r User) UpdateUser() error {
	err := db.C(USER).Update(bson.M{"_id": r.ID}, &r)
	return err
}

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
	ITEM = config.Item
	USER = config.User
	DB = config.Database
	Connect(url)
}
