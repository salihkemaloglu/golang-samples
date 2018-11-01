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

//DB string
var DB string

// FindAll list of Items
func (r Item) FindAll() ([]byte, error) {
	var items []Item
	err := db.C("item").Find(bson.M{}).All(&items)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(items)
	return data, err
}

// FindAll list of users
func (r User) FindAll() ([]byte, error) {
	var user []User
	err := db.C("user").Find(bson.M{}).All(&user)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(user)
	return data, err
}

// FindById a Items by its id
func (r Item) FindById() ([]byte, error) {
	err := db.C("item").FindId(bson.ObjectIdHex(r.ItemId)).One(&r)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(r)
	return data, err
}

// FindById a users by its id
func (r User) FindById() ([]byte, error) {
	err := db.C("user").FindId(bson.ObjectIdHex(*r.Password)).One(&r)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(r)
	return data, err
}

// Insert a Items into database
func (r Item) Insert() error {
	r.ID = bson.NewObjectId()
	err := db.C("item").Insert(&r)
	return err
}

// Insert a users into database
func (r User) Insert() error {
	r.ID = bson.NewObjectId()
	err := db.C("user").Insert(&r)
	return err
}

// Delete an existing Items
func (r Item) Delete() error {
	err := db.C("item").Remove(&r)
	return err
}

// Delete an existing users
func (r User) Delete() error {
	err := db.C("user").Remove(&r)
	return err
}

// Update an existing Items
func (r Item) Update() error {
	err := db.C("item").Update(bson.M{"_id": r.ID}, &r)
	return err
}

// Update an existing users
func (r User) Update() error {
	err := db.C("user").Update(bson.M{"_id": r.ID}, &r)
	return err
}

// Login a user
func (r User) Login() ([]byte, error) {
	err := db.C("user").Find(bson.M{"username": r.Username, "password": r.Password}).One(&r)
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(r)
	return data, err
}

//Register user
func (r User) Register() error {
	r.ID = bson.NewObjectId()
	err := db.C("user").Insert(&r)
	return err
}

//CheckUser user login
func (r User) CheckUser() bool {
	err := db.C("user").Find(bson.M{"username": r.Username}).One(&r)
	if err != nil {
		return true
	}
	return false
}

//Connect Establish a connection to database
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

//LoadConfiguration Parse the configuration file 'config.toml', and establish a connection to DB
func LoadConfiguration() {
	var url = os.Getenv("HOST_ENV")
	DB = os.Getenv("DATABASE_ENV")
	Connect(url)
}
