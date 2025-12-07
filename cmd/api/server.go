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
	utils "github.com/GuiFernandess7/risa/pkg/utils"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		fmt.Println("Database Error: ", err)
		return
	}

	e := echo.New()
	e.Validator = &utils.CustomValidator{Validator: validator.New()}
	e = middlewares.ApplySecurityMiddlewares(e)

	router.InitRoutes(db, e)
	srv := &http.Server{
		Addr:         ":" + port,
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
