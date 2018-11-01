package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

//CreateTokenEndpoint user token creation
func CreateTokenEndpoint(user User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"username": user.Username,
		"password": user.Password,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return ""
	}
	return tokenString
}

//ValidateMiddleware token validation
func ValidateMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authorizationHeader := req.Header.Get("authorization")
		if authorizationHeader != "" {
			bearerToken := strings.Split(authorizationHeader, " ")
			if len(bearerToken) == 2 {
				token, error := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("There was an error")
					}
					return []byte("secret"), nil
				})
				if error != nil {
					json.NewEncoder(w).Encode(Exception{Message: error.Error()})
					return
				}
				if token.Valid {
					context.Set(req, "decoded", token.Claims)
					next(w, req)
				} else {
					json.NewEncoder(w).Encode(Exception{Message: "Invalid authorization token"})
				}
			}
		} else {
			json.NewEncoder(w).Encode(Exception{Message: "An authorization header is required"})
		}
	})
}

//UserFieldValidation user struct validation
func UserFieldValidation(user User) string {
	if user.Username == nil {
		return "Username field is missing"
	} else if user.Password == nil {
		return "Password field is missing"
	} else if strings.TrimSpace(*user.Username) == "" {
		return "Username can not be empty"
	} else if strings.TrimSpace(*user.Password) == "" {
		return "Password can not be empty"
	} else {
		return "ok"
	}
}

//ItemFildValidation item struct validaiton
func ItemFildValidation(item Item) string {
	if item.Name == nil {
		return "Name fields is missing"
	} else if item.Value == nil {
		return "Value files is missing"
	} else if strings.TrimSpace(*item.Name) == "" {
		return "Name can not be empty"
	} else if strings.TrimSpace(*item.Value) == "" {
		return "Value can not be empty"
	} else if len(item.Description) > 25 {
		return "Description length can not be bigger than 25"
	} else {
		return "ok"
	}
}
