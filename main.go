package main

import (
	"log"
	"messaging-service/controllers"
	"messaging-service/database"
	"messaging-service/websockets"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.InitDB()
	websockets.InitWSManager()

	router := mux.NewRouter()

	router.HandleFunc("/register", controllers.Register)
	router.HandleFunc("/login", controllers.Login)
	router.HandleFunc("/users", controllers.ListUsers)
	router.HandleFunc("/messages", controllers.GetMessageHistory).Methods("GET")

	router.HandleFunc("/ws", websockets.HandleConnections)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Allow all origins
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Wrap the router with the CORS middleware
	handler := c.Handler(router)

	// Start the server
	log.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
