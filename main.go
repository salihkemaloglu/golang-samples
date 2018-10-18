package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gopkg.in/mgo.v2/bson"
)

func UserRegister(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	} else if check := UserFieldValidation(user); check != "ok" {
		respondWithError(w, http.StatusBadRequest, check)
		return
	}
	if err := CheckUser(user); err != true {
		respondWithError(w, http.StatusInternalServerError, "There is alreasy an account for: "+*user.Username)
		return
	}
	user.ID = bson.NewObjectId()
	if err := Register(user); err != true {
		respondWithError(w, http.StatusInternalServerError, "Registration problem!")
		return
	}
	user.Token = ""
	respondWithJson(w, http.StatusCreated, user)
}

func UserLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	} else if check := UserFieldValidation(user); check != "ok" {
		respondWithError(w, http.StatusBadRequest, check)
		return
	}
	user, err := Login(user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid username or password")
		return
	}
	token := CreateTokenEndpoint(user)
	if token == "" {
		user.Token = "Token could not created"
	}
	user.Token = token
	*user.Username = ""
	*user.Password = ""
	respondWithJson(w, http.StatusOK, user)
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
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	repo := User{}
	var users, err = FindAllUser(repo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, users)
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
	} else if check := ItemFildValidation(item); check != "ok" {
		respondWithError(w, http.StatusBadRequest, check)
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
	defer r.Body.Close()
	var itemGet Item
	if err := json.NewDecoder(r.Body).Decode(&itemGet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	} else if check := ItemFildValidation(itemGet); check != "ok" {
		respondWithError(w, http.StatusBadRequest, check)
		return
	}
	params := mux.Vars(r)
	repo := Item{ItemId: params["id"]}
	item, err := FindById(repo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	itemGet.ID = item.ID
	if err := Update(itemGet); err != nil {
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
	myRouter.HandleFunc("/login", UserLogin).Methods("POST")
	myRouter.HandleFunc("/register", UserRegister).Methods("POST")
	myRouter.HandleFunc("/item", ValidateMiddleware(GetAll)).Methods("GET")
	myRouter.HandleFunc("/user", ValidateMiddleware(GetAllUser)).Methods("GET")
	myRouter.HandleFunc("/item/{id}", ValidateMiddleware(GetById)).Methods("GET")
	myRouter.HandleFunc("/item", ValidateMiddleware(InsertItem)).Methods("POST")
	myRouter.HandleFunc("/item/{id}", ValidateMiddleware(UpdateItem)).Methods("PUT")
	myRouter.HandleFunc("/item/{id}", ValidateMiddleware(DeleteItem)).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", cors.AllowAll().Handler(myRouter)))
}

func main() {
	LoadConfiguration()
	handleRequests()
}
