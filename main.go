package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Items []Item

var db *mgo.Database

var COLLECTION string
var DB string

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!Service is working")
}

// GET list of items
func GetAll(w http.ResponseWriter, r *http.Request) {
	repo := Item{}
	var items, err = FindAll(repo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, items)
}

// GET a item by its ID
func GetById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repo := Item{ItemId: params["id"]}
	item, err := FindById(repo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	respondWithJson(w, http.StatusOK, item)
}

// POST a new item
func InsertItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	item.ID = bson.NewObjectId()
	if err := Insert(item); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, item)
}

// PUT update an existing item
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repo := Item{ItemId: params["id"]}
	item, err := FindById(repo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	var bsonid = item.ID
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	item.ID = bsonid
	if err := Update(item); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

// DELETE an existing item
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repo := Item{ItemId: params["id"]}
	item, err := FindById(repo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	if err := Delete(item); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	respondWithJson(w, code, map[string]string{"error": msg})
}

func respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/item", GetAll).Methods("GET")
	myRouter.HandleFunc("/item/{id}", GetById).Methods("GET")
	myRouter.HandleFunc("/item", InsertItem).Methods("POST")
	myRouter.HandleFunc("/item/{id}", UpdateItem).Methods("PUT")
	myRouter.HandleFunc("/item/{id}", DeleteItem).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":3001", myRouter))
}

func main() {
	LoadConfiguration()
	handleRequests()
}
