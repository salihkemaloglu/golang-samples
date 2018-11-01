package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"

	"github.com/salihkemaloglu/GitRepo"

	"github.com/rs/cors"
	"goji.io"
	"goji.io/pat"
)

//UserRegister User registration
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
	var e UserRepository = user
	if err := e.CheckUser(); err != true {
		respondWithError(w, http.StatusInternalServerError, "There is alreasy an account for: "+*user.Username)
		return
	}
	if err := e.Register(); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Registration problem!"+err.Error())
		return
	}
	user.Token = ""
	respondWithJson(w, http.StatusCreated, user)
}

//UserLogin user login
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
	var e UserRepository = user
	data, err := e.Login()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid username or password")
		return
	}
	err = json.Unmarshal(data, &user)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
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

//GetAll  GET list of items
func GetAll(w http.ResponseWriter, r *http.Request) {
	repo := Item{}
	var e BaseRepository = repo
	var items, err = e.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var item []Item
	err = json.Unmarshal(items, &item)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, item)
}

//GetAllUser  get all user list
func GetAllUser(w http.ResponseWriter, r *http.Request) {
	repo := User{}
	var e BaseRepository = repo
	var users, err = e.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	var user []User
	err = json.Unmarshal(users, &user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, user)
}

//GetById GET a item by its ID
func GetById(w http.ResponseWriter, r *http.Request) {
	params := pat.Param(r, "id")
	repo := Item{ItemId: params}
	var e BaseRepository = repo
	item, err := e.FindById()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	err = json.Unmarshal(item, &repo)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, repo)
}

//InsertItem  POST a new item
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
	var e BaseRepository = item
	if err := e.Insert(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, item)
}

//UpdateItem  PUT update an existing item
func UpdateItem(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var itemGet Item
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&itemGet); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	} else if check := ItemFildValidation(itemGet); check != "ok" {
		respondWithError(w, http.StatusBadRequest, check)
		return
	}

	params := pat.Param(r, "id")
	repo := Item{ItemId: params}
	var e BaseRepository = repo
	data, err := e.FindById()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	er := json.Unmarshal(data, &item)
	if er != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid username or password")
		return
	}
	itemGet.ID = item.ID
	e = itemGet
	if err := e.Update(); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, map[string]string{"result": "success"})
}

//DeleteItem DELETE an existing item
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	params := pat.Param(r, "id")
	repo := Item{ItemId: params}
	var item Item
	var e BaseRepository = repo
	data, err := e.FindById()
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item ID")
		return
	}
	er := json.Unmarshal(data, &item)
	if er != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid username or password")
		return
	}
	e = item
	if err := e.Delete(); err != nil {
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
	mux := goji.NewMux()
	mux.HandleFunc(pat.Post("/login"), UserLogin)
	mux.HandleFunc(pat.Post("/register"), UserRegister)
	mux.HandleFunc(pat.Get("/item"), ValidateMiddleware(GetAll))
	mux.HandleFunc(pat.Get("/user"), ValidateMiddleware(GetAllUser))
	mux.HandleFunc(pat.Get("/item/:id"), ValidateMiddleware(GetById))
	mux.HandleFunc(pat.Post("/item"), ValidateMiddleware(InsertItem))
	mux.HandleFunc(pat.Put("/item/:id"), ValidateMiddleware(UpdateItem))
	mux.HandleFunc(pat.Delete("/item/:id"), ValidateMiddleware(DeleteItem))
	log.Fatal(http.ListenAndServe(":8080", cors.AllowAll().Handler(mux)))
}

func main() {
	fmt.Println(mdata.TestMe())
	LoadConfiguration()
	handleRequests()

}
