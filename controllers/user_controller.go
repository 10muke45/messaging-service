package controllers

import (
	"encoding/json"
	"messaging-service/services"
	"messaging-service/utils"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&reqBody)

	err := services.CreateUser(reqBody.Username, reqBody.Password)
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&reqBody)

	if !services.ValidateUser(reqBody.Username, reqBody.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(reqBody.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func ListUsers(w http.ResponseWriter, r *http.Request) {
	res, err := services.ListUsers()
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string][]string{"data": res})
}

func GetMessageHistory(w http.ResponseWriter, r *http.Request) {
	// Get the token from the request header
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}
	senderUserName, ok := utils.ValidateToken(tokenString)
	if !ok {
		http.Error(w, "Authorization token required", http.StatusUnauthorized)
		return
	}
	receiverUserName := r.URL.Query().Get("username")
	messages, err := services.GetMessageHistoryForUser(senderUserName, receiverUserName)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(messages)
}
