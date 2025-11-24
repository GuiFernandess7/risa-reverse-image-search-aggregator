package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	router "github.com/GuiFernandess7/risa/internal"
	middlewares "github.com/GuiFernandess7/risa/internal/middlewares"
	sqlconnect "github.com/GuiFernandess7/risa/internal/repository/database"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	port := os.Getenv("API_PORT")

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Database Error: ", err)
		return
	}

	e := echo.New()
	e = middlewares.ApplyMiddlewares(e)

	router.InitRoutes(db, e)
	srv := &http.Server{
		Addr:         port,
		Handler:      e,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Println("Server Listening on port:", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalln("error starting server:", err)
	}
}
